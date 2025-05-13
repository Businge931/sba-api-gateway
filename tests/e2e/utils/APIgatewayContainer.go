package utils

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Helper function to set up the API Gateway container from the Dockerfile
func SetupAPIGatewayContainer(ctx context.Context, t *testing.T, networkName string, authServiceC, oddsServiceC testcontainers.Container) (testcontainers.Container, error) {
	t.Log("Building and setting up API Gateway container from Dockerfile")

	// Get the api-gateway project directory
	apiGatewayDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		return nil, fmt.Errorf("failed to get api gateway directory: %w", err)
	}

	// Get actual container IPs within Docker network instead of relying on DNS
	authIP, err := authServiceC.ContainerIP(ctx)
	if err != nil {
		t.Logf("Error getting auth service container IP: %v, falling back to container name", err)
		authIP = "auth-service" // Fallback to container network alias
	}

	oddsIP, err := oddsServiceC.ContainerIP(ctx)
	if err != nil {
		t.Logf("Error getting odds service container IP: %v, falling back to container name", err)
		oddsIP = "odds-service" // Fallback to container network alias
	}

	// Use IP addresses directly for services since DNS resolution is problematic
	authServiceAddr := fmt.Sprintf("%s:50051", authIP)
	oddsServiceAddr := fmt.Sprintf("%s:50052", oddsIP)

	t.Logf("Using auth service at container IP: %s", authServiceAddr)
	t.Logf("Using odds service at container IP: %s", oddsServiceAddr)

	// API Gateway container configuration
	req := testcontainers.ContainerRequest{
		Name: "api-gateway",
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    apiGatewayDir,
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"8080/tcp"},
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"api-gateway"},
		},
		// Increase timeout for startup to ensure services are ready
		WaitingFor: wait.ForHTTP("/health").WithPort("8080/tcp").WithStartupTimeout(3 * time.Minute),
		Env: map[string]string{
			// Explicitly specify service addresses with format that works in Docker
			"AUTH_SERVICE_ADDR": authServiceAddr,
			"ODDS_SERVICE_ADDR": oddsServiceAddr,
			
			// Alternative formats for service addresses (in case app uses different env vars)
			"AUTH_SERVICE_HOST": authIP,
			"AUTH_SERVICE_PORT": "50051",
			"ODDS_SERVICE_HOST": oddsIP,
			"ODDS_SERVICE_PORT": "50052",
			
			// Explicit database connection info for services
			"DB_HOST": "postgres",
			"DB_PORT": "5432",
			"DB_USER": "postgres",
			"DB_PASSWORD": "postgres",
			"AUTH_DB_NAME": "sba_user_accounts",
			"ODDS_DB_NAME": "sba_odds",

			// JWT settings
			"JWT_SECRET": "testsecretkey",
			"JWT_EXPIRY": "24h",

			// Force using DNS resolver for gRPC
			"GRPC_DNS_RESOLVER": "native",
			"GRPC_GO_REQUIRE_AUTHORITY": "false", // Try to bypass authority issues with gRPC
			"GRPC_GO_RETRY_INDEFINITELY": "true",  // Keep trying to connect
			"GRPC_RETRY_TIMEOUT": "10s",           // Give more time for retries

			// Enhanced debugging and logging
			"LOG_LEVEL": "debug",
			"DEBUG": "true",
			"TRACE": "true",
			"ENABLE_REQUEST_LOGGING": "true",
			"ENABLE_ERROR_DETAILS": "true",
			"GRPC_GO_LOG_VERBOSITY_LEVEL": "99",
			"GRPC_GO_LOG_SEVERITY_LEVEL": "info",
			
			// Environment mode
			"ENV": "test",
			"NODE_ENV": "test",
			"GO_ENV": "test",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start API gateway container: %w", err)
	}

	// Get container connection info
	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())
	t.Logf("API Gateway container is ready at: %s", endpoint)

	// Capture logs to help with debugging
	time.Sleep(1 * time.Second) // Give service time to log any startup issues
	gatewayLogs, gatewayLogsErr := container.Logs(ctx)
	if gatewayLogsErr == nil {
		gatewayLogContent, _ := io.ReadAll(gatewayLogs)
		t.Logf("Initial API Gateway logs: %s", gatewayLogContent)
	} else {
		t.Logf("Error getting API Gateway logs: %v", gatewayLogsErr)
	}

	// Verify API Gateway is healthy
	t.Log("Checking API Gateway health status...")
	health := VerifyServiceHealth(endpoint + "/health", t)
	if !health {
		t.Log("Warning: API Gateway health check failed but continuing with test")
	}
	return container, nil
}