package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Businge931/sba-api-gateway/internal/api/adapters"
	"github.com/Businge931/sba-api-gateway/internal/api/handlers"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
)

func SetupRoutes(authService service.AuthService, oddsService service.OddsService) *mux.Router {
	router := mux.NewRouter()

	// Create adapters for the services
	authAdapter := adapters.NewAuthServiceAdapter(authService)
	oddsAdapter := adapters.NewOddsServiceAdapter(oddsService)

	// Register handlers
	authHandler := handlers.NewAuthHandler(authAdapter)
	oddsHandler := handlers.NewOddsHandler(oddsAdapter)

	// Public endpoints
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")

	// Health check endpoint for container testing
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Protected endpoints (require token verification)
	// Odds management endpoints
	router.Handle("/api/odds/create", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.CreateOdds))).Methods("POST")
	router.Handle("/api/odds/read", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.ReadOdds))).Methods("GET")
	router.Handle("/api/odds/update", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.UpdateOdds))).Methods("PUT")
	router.Handle("/api/odds/delete", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.DeleteOdds))).Methods("DELETE")

	// Legacy endpoints for backward compatibility
	router.Handle("/create", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.CreateOdds))).Methods("POST")
	router.Handle("/read", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.ReadOdds))).Methods("GET")
	router.Handle("/update", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.UpdateOdds))).Methods("PUT")
	router.Handle("/delete", authHandler.VerifyTokenMiddleware(http.HandlerFunc(oddsHandler.DeleteOdds))).Methods("DELETE")

	return router
}
