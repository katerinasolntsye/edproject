package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// NewCORS creates a new CORS middleware with appropriate settings
func NewCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodHead,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:          300, // Maximum age for preflight requests (5 minutes)
	})
}
