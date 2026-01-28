package app

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/katerinasolntsye/fulleng/internal/config"
	"github.com/katerinasolntsye/fulleng/internal/handler"
	"github.com/katerinasolntsye/fulleng/internal/middleware"
	"github.com/katerinasolntsye/fulleng/internal/repository/postgres"
	"github.com/katerinasolntsye/fulleng/internal/service"
)

type App struct {
	config *config.Config
	router *mux.Router
	conn   *pgx.Conn
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize() error {
	// Load configuration
	a.config = config.Load()

	// Connect to database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, a.config.Database.URL)
	if err != nil {
		return err
	}
	a.conn = conn

	// Initialize layers
	repo := postgres.NewRepository(a.conn)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	// Setup routes
	a.router = mux.NewRouter()
	a.setupRoutes(h)

	return nil
}

func (a *App) setupRoutes(h *handler.Handler) {
	api := a.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/incoming", h.GetIncomingPostback).Methods("GET")
	api.HandleFunc("/tracker", h.GetTracker).Methods("GET")
	api.HandleFunc("/sendpostback", h.GetSendPostback).Methods("GET")
	api.HandleFunc("/signup", h.Signup).Methods("POST")
	api.HandleFunc("/signin", h.Signin).Methods("POST")
	api.HandleFunc("/user/{id}", h.GetUser).Methods("GET")
	api.HandleFunc("/user/{id}", h.UpdateUser).Methods("POST")
	api.HandleFunc("/user/{id}/creds", h.UpdateUserCredentials).Methods("POST")
}

func (a *App) Run() error {
	defer a.conn.Close(context.Background())

	// Apply CORS middleware
	corsMiddleware := middleware.NewCORS()
	handler := corsMiddleware.Handler(a.router)

	log.Printf("Server starting on port %s", a.config.Server.Port)
	return http.ListenAndServe(a.config.Server.Port, handler)
}
