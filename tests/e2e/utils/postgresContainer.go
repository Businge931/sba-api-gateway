package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupPostgresContainer(ctx context.Context, t *testing.T, networkName string) (testcontainers.Container, error) {
	t.Log("Setting up PostgreSQL container")

	// PostgreSQL container configuration
	req := testcontainers.ContainerRequest{
		Name:         "postgres-test",
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Networks:     []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"postgres"},
		},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start PostgreSQL container: %w", err)
	}

	// Get container connection info
	ip, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	endpoint := fmt.Sprintf("%s:%s", ip, mappedPort.Port())
	t.Logf("PostgreSQL container is ready at: %s", endpoint)

	// Give PostgreSQL time to fully initialize before running our SQL commands
	t.Log("Waiting for PostgreSQL to fully initialize...")
	time.Sleep(2 * time.Second)

	// Initialize the required databases by running SQL commands in the container
	err = CreateAndInitDatabases(ctx, t, container)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize databases: %w", err)
	}

	// Verify databases were created successfully
	t.Log("Verifying database creation...")
	if err = VerifyDatabasesExist(ctx, t, container); err != nil {
		return nil, fmt.Errorf("database verification failed: %w", err)
	}

	return container, nil
}
