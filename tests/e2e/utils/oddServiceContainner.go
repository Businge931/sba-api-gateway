package utils

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Helper function to set up the Odds Service container from the Dockerfile
func SetupOddsServiceContainer(ctx context.Context, t *testing.T, networkName string) (testcontainers.Container, error) {
	t.Log("Building and setting up Odds Service container from Dockerfile")

	// Get the parent directory of the api-gateway project to find crud-ops service
	apiGatewayDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		return nil, fmt.Errorf("failed to get api gateway directory: %w", err)
	}

	// Path to crud-ops service
	crudOpsDir := filepath.Join(filepath.Dir(apiGatewayDir), "sba-crud-ops", "sba-crud-ops")
	t.Logf("Using crud ops path: %s", crudOpsDir)

	// Container configuration
	req := testcontainers.ContainerRequest{
		Name: "odds-service",
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    crudOpsDir,
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"50052/tcp"},
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"odds-service"},
		},
		Env: map[string]string{
			// The odds service uses a single DB_ADDR connection string rather than individual variables
			"DB_ADDR": "postgresql://postgres:postgres@postgres:5432/sba_odds?sslmode=disable",
			
			// Add individual DB parameters as well in case the service falls back to these
			"DB_HOST":     "postgres",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASSWORD": "postgres",
			"DB_NAME":     "sba_odds",
			
			// Enhanced logging for debugging
			"LOG_LEVEL":    "debug",
			"TRACE_SQL":    "true",    // Enable SQL query tracing
			"ENABLE_TRACE": "true",    // Enable general tracing
			"LOG_FORMAT":   "json",     // Use structured JSON logging
			
			// Force using native DNS resolver for gRPC
			"GRPC_DNS_RESOLVER":         "native",
			"GRPC_GO_LOG_VERBOSITY_LEVEL": "99",
			"GRPC_GO_LOG_SEVERITY_LEVEL":  "info",
			
			// Database initialization flags
			"AUTO_MIGRATE":  "true",    // Auto migrate database schema
			"INIT_DB_TABLES": "true",   // Force table initialization
			"SEED_TEST_DATA": "true",   // Seed with test data for testing
		},
		WaitingFor: wait.ForLog("Starting odds-service gRPC server on port 50052").WithStartupTimeout(2 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start odds service container: %w", err)
	}

	// Get container connection info
	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "50052/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	endpoint := fmt.Sprintf("%s:%s", ip, mappedPort.Port())
	t.Logf("Odds service container is ready at: %s", endpoint)

	// Capture logs to help with debugging
	time.Sleep(1 * time.Second) // Give service time to log any startup issues
	oddsLogs, oddsLogsErr := container.Logs(ctx)
	if oddsLogsErr == nil {
		oddsLogContent, _ := io.ReadAll(oddsLogs)
		t.Logf("Odds service logs: %s", oddsLogContent)
	} else {
		t.Logf("Error getting odds service logs: %v", oddsLogsErr)
	}
	return container, nil
}
