package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Businge931/sba-api-gateway/internal/api"
	"github.com/Businge931/sba-api-gateway/internal/api/middleware"
	"github.com/Businge931/sba-api-gateway/internal/app/service"

	server "github.com/Businge931/sba-api-gateway/internal/server/grpc"
)

// getEnv retrieves the value of the environment variable named by the key
// If the variable is not present, returns the fallback value
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func main() {
	// Set up logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Get service hosts and ports from environment variables or use defaults
	oddsServiceHost := getEnv("ODDS_SERVICE_HOST", "localhost")
	oddsServicePort := getEnv("ODDS_SERVICE_PORT", "50052")
	authServiceHost := getEnv("AUTH_SERVICE_HOST", "localhost")
	authServicePort := getEnv("AUTH_SERVICE_PORT", "50051")
	
	oddsAddr := oddsServiceHost + ":" + oddsServicePort
	authAddr := authServiceHost + ":" + authServicePort
	
	// Connect to odds service
	log.Infof("Connecting to odds service at %s", oddsAddr)
	oddsConn, err := grpc.Dial(oddsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to odds service: %v", err)
	}
	defer oddsConn.Close()
	log.Info("Connected to odds service")

	// Connect to auth service
	log.Infof("Connecting to auth service at %s", authAddr)
	authConn, err := grpc.Dial(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()
	log.Info("Connected to auth service")

	// Initialize gRPC clients
	authServer := server.NewAuthServer(authConn)
	oddsServer := server.NewOddsServer(oddsConn)

	// Initialize services
	authService := service.NewAuthService(authServer)
	oddsService := service.NewOddsService(oddsServer)

	// Setup routes
	router := api.SetupRoutes(authService, oddsService)

	// Register middleware
	router.Use(middleware.LoggingMiddleware)

	// Add CORS middleware
	corsHandler := middleware.GetCORSConfig()
	handler := corsHandler.Handler(router)

	// Start the server
	api.Start(handler, ":8080")
}
