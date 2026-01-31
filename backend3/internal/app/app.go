package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"github.com/katerinasolntsye/fulleng/internal/config"
	"github.com/katerinasolntsye/fulleng/internal/handler"
	"github.com/katerinasolntsye/fulleng/internal/middleware"
	"github.com/katerinasolntsye/fulleng/internal/repository/sqlite"
	"github.com/katerinasolntsye/fulleng/internal/service"
)

type App struct {
	config     *config.Config
	router     *mux.Router
	db         *sql.DB
	jwtService *service.JWTService
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize() error {
	// Load configuration
	a.config = config.Load()

	// Connect to SQLite database
	db, err := sql.Open("sqlite3", a.config.Database.Path)
	if err != nil {
		return err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return err
	}

	a.db = db

	// Initialize JWT service
	a.jwtService = service.NewJWTService(
		a.config.JWT.Secret,
		a.config.JWT.AccessExpiration,
		a.config.JWT.RefreshExpiration,
	)

	// Initialize layers
	repo := sqlite.NewRepository(a.db)
	svc := service.NewService(repo, a.jwtService)
	h := handler.NewHandler(svc)

	// Setup routes
	a.router = mux.NewRouter()
	a.setupRoutes(h)

	return nil
}

func (a *App) setupRoutes(h *handler.Handler) {
	api := a.router.PathPrefix("/api/v1").Subrouter()

	// Публичные роуты (без авторизации)
	api.HandleFunc("/signup", h.Signup).Methods("POST")
	api.HandleFunc("/signin", h.Signin).Methods("POST")
	api.HandleFunc("/refresh", h.RefreshToken).Methods("POST")

	// Защищенные роуты (с авторизацией)
	authMiddleware := middleware.NewAuthMiddleware(a.jwtService)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMiddleware.RequireAuth)

	protected.HandleFunc("/logout", h.Logout).Methods("POST")
	protected.HandleFunc("/user/{id}", h.GetUser).Methods("GET")
	protected.HandleFunc("/user/{id}", h.UpdateUser).Methods("POST")
	protected.HandleFunc("/user/{id}/creds", h.UpdateUserCredentials).Methods("POST")
	protected.HandleFunc("/incoming", h.GetIncomingPostback).Methods("GET")
	protected.HandleFunc("/tracker", h.GetTracker).Methods("GET")
	protected.HandleFunc("/sendpostback", h.GetSendPostback).Methods("GET")
}

func (a *App) Run() error {
	defer a.db.Close()

	// Apply CORS middleware
	corsMiddleware := middleware.NewCORS()
	handler := corsMiddleware.Handler(a.router)

	log.Printf("Server starting on port %s", a.config.Server.Port)
	log.Printf("Using SQLite database: %s", a.config.Database.Path)

	return http.ListenAndServe(a.config.Server.Port, handler)
}
