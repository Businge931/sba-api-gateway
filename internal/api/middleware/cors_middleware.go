package middleware

import (
	"github.com/rs/cors"
)

// GetCORSConfig returns a CORS handler with predefined options
func GetCORSConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins including Postman
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "Origin", "User-Agent"},
		AllowCredentials: true,
		Debug:           true, // Enable debug logging for CORS issues
	})
}
