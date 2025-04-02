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

func main() {
	// Set up logging
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Connect to gRPC services
	// Connect to odds service on port 50052
	oddsConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to odds service: %v", err)
	}
	defer oddsConn.Close()
	log.Info("Connected to odds service")

	// Connect to auth service on port 50051
	authConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
