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

// Helper function to set up the Auth Service container from the Dockerfile
func SetupAuthServiceContainer(ctx context.Context, t *testing.T, networkName string) (testcontainers.Container, error) {
	t.Log("Building and setting up Auth Service container from Dockerfile")

	// Get the parent directory of the api-gateway project to find user-accounts service
	apiGatewayDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		return nil, fmt.Errorf("failed to get api gateway directory: %w", err)
	}

	// Path to user-accounts service
	userAccountsDir := filepath.Join(filepath.Dir(apiGatewayDir), "sba-user-accounts", "sba-user-accounts")
	t.Logf("Using user accounts path: %s", userAccountsDir)

	// Container configuration
	req := testcontainers.ContainerRequest{
		Name: "auth-service",
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    userAccountsDir,
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"50051/tcp"},
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"auth-service"},
		},
		Env: map[string]string{
			// The auth service uses separate environment variables
			"DB_HOST":     "postgres",
			"DB_PORT":     "5432",
			"DB_USER":     "postgres",
			"DB_PASSWORD": "postgres",
			"DB_NAME":     "sba_user_accounts",
			
			// JWT configuration settings
			"JWT_SECRET":  "testsecretkey",
			"JWT_EXPIRY":  "24h",
			
			// Additional logging for debugging
			"LOG_LEVEL":    "debug",
			"TRACE_LEVEL":  "sql",
			"ENABLE_TRACE": "true",
			
			// Additional flags for troubleshooting
			"SKIP_PASSWORD_VALIDATION": "true", // Make password validation less strict for tests
			"INIT_DB_TABLES": "true",         // Force table initialization
		},
		WaitingFor: wait.ForLog("Starting gRPC server on port 50051").WithStartupTimeout(2 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start auth service container: %w", err)
	}

	// Get container connection info
	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "50051/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	endpoint := fmt.Sprintf("%s:%s", ip, mappedPort.Port())
	t.Logf("Auth service container is ready at: %v", endpoint)

	// Capture logs to help with debugging
	time.Sleep(1 * time.Second) // Give service time to log any startup issues
	authLogs, authLogsErr := container.Logs(ctx)
	if authLogsErr == nil {
		authLogContent, _ := io.ReadAll(authLogs)
		t.Logf("Auth service logs: %s", authLogContent)
	} else {
		t.Logf("Error getting auth service logs: %v", authLogsErr)
	}

	return container, nil
}
