package e2e

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/Businge931/sba-api-gateway/tests/e2e/utils"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/network"
)

// TestWithContainers demonstrates how to use Testcontainers to test the API gateway
// with containerized microservices
func TestWithContainers(t *testing.T) {
	// Skip this test when running unit tests
	// This test is meant to be run in a separate integration test pipeline
	if os.Getenv("SKIP_CONTAINER_TESTS") == "true" {
		t.Skip("Skipping container-based test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create a Docker network for the tests
	t.Log("Creating Docker network for test containers")
	networkName := "sba-test-network"
	
	// Create a new Docker network
	dockerNetwork, err := network.New(ctx)
	if err != nil {
		t.Fatalf("Failed to create Docker network: %s", err)
	}
	// Get the network name for container connection
	networkName = dockerNetwork.Name
	
	defer func() {
		if err := dockerNetwork.Remove(ctx); err != nil {
			t.Logf("Failed to remove Docker network: %s", err)
		}
	}()

	// Start PostgreSQL container for test databases
	t.Log("Setting up PostgreSQL database container")
	postgresC, err := utils.SetupPostgresContainer(ctx, t, networkName)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %s", err)
	}
	defer func() {
		if err := postgresC.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate PostgreSQL container: %s", err)
		}
	}()

	// Wait a bit for the database to be fully initialized and tables created
	t.Log("Waiting for database initialization to complete...")
	time.Sleep(3 * time.Second)

	// Create the necessary databases manually once the container is up
	t.Log("Creating necessary databases manually")

	// Create sba_user_accounts database
	_, _, err = postgresC.Exec(ctx, []string{"psql", "-U", "postgres", "-c", "CREATE DATABASE sba_user_accounts;"})
	if err != nil {
		t.Logf("Note: Database creation error (may already exist): %v", err)
	}

	// Create admin role
	_, _, err = postgresC.Exec(ctx, []string{"psql", "-U", "postgres", "-c", "CREATE ROLE admin WITH LOGIN PASSWORD 'admin';"})
	if err != nil {
		t.Logf("Note: Role creation error (may already exist): %v", err)
	}

	// Create sba_crud_ops database
	_, _, err = postgresC.Exec(ctx, []string{"psql", "-U", "postgres", "-c", "CREATE DATABASE sba_crud_ops;"})
	if err != nil {
		t.Logf("Note: Database creation error (may already exist): %v", err)
	}

	// Grant privileges
	_, _, err = postgresC.Exec(ctx, []string{"psql", "-U", "postgres", "-c", "GRANT ALL PRIVILEGES ON DATABASE sba_crud_ops TO admin;"})
	if err != nil {
		t.Logf("Note: Grant privileges error: %v", err)
	}

	// Additional sleep to ensure database is fully ready
	time.Sleep(2 * time.Second)

	// Create users table in auth database with the schema expected by the Auth service
	t.Log("Creating auth users table...")
	_, _, err = postgresC.Exec(ctx, []string{
		"psql", "-U", "postgres", "-d", "sba_user_accounts", "-c",
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			hashed_password VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			is_email_verified BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`})
	if err != nil {
		t.Logf("Error creating auth users table: %v", err)
	}

	// Add test user for login tests with plaintext password for simplicity
	// IMPORTANT: The username must be in email format as the Auth service validates it
	// Hash for 'password123' generated with bcrypt cost 10 which is the standard default
	t.Log("Adding test user to auth database...")
	_, _, err = postgresC.Exec(ctx, []string{
		"psql", "-U", "postgres", "-d", "sba_user_accounts", "-c",
		`INSERT INTO users (username, hashed_password, email, first_name, last_name, is_email_verified, updated_at) 
		 VALUES ('test@example.com', '$2a$10$iil1PbnIHBsqv4yRemxcg.qTXY/78UvI/BmHRPz4MuJWoAlx3VRtW', 'test@example.com', 'Test', 'User', TRUE, now()) 
		 ON CONFLICT (username) DO NOTHING;`})
	if err != nil {
		t.Logf("Error adding test user: %v", err)
	}

	// Set up database schema for odds service
	t.Log("Setting up schema for sba_crud_ops database...")
	_, _, err = postgresC.Exec(ctx, []string{
		"psql", "-U", "postgres", "-d", "sba_crud_ops", "-c",
		`CREATE TABLE IF NOT EXISTS odds (
			id SERIAL PRIMARY KEY,
			league VARCHAR(255) NOT NULL,
			home_team VARCHAR(255) NOT NULL,
			away_team VARCHAR(255) NOT NULL,
			home_team_win_odds DECIMAL(5,2) NOT NULL,
			away_team_win_odds DECIMAL(5,2) NOT NULL,
			draw_odds DECIMAL(5,2) NOT NULL,
			game_date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`})
	if err != nil {
		t.Logf("Error creating odds table: %v", err)
	}

	t.Log("Database schema setup completed")

	// Within the Docker network, containers can access the database using its network alias 'postgres'
	t.Log("Using 'postgres' as hostname for database connection within Docker network")

	// Get postgres container logs to verify initialization
	dbLogs, logsErr := postgresC.Logs(ctx)
	if logsErr == nil {
		logContent, _ := io.ReadAll(dbLogs)
		t.Logf("Postgres initialization logs: %s", logContent)
	}

	// Start Auth Service container on the same network
	authServiceC, err := utils.SetupAuthServiceContainer(ctx, t, networkName)
	if err != nil {
		t.Fatalf("Failed to start auth service container: %s", err)
	}
	defer func() {
		if err := authServiceC.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate auth service container: %s", err)
		}
	}()

	// Start Odds Service container on the same network
	oddsServiceC, err := utils.SetupOddsServiceContainer(ctx, t, networkName)
	if err != nil {
		t.Fatalf("Failed to start odds service container: %s", err)
	}
	defer func() {
		if err := oddsServiceC.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate odds service container: %s", err)
		}
	}()

	// Start API Gateway container with environment variables pointing to the other services
	apiGatewayC, err := utils.SetupAPIGatewayContainer(ctx, t, networkName, authServiceC, oddsServiceC)
	if err != nil {
		t.Fatalf("Failed to start API gateway container: %s", err)
	}
	defer func() {
		if err := apiGatewayC.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate API gateway container: %s", err)
		}
	}()

	// Get API Gateway endpoint for testing
	apiEndpoint, err := apiGatewayC.Endpoint(ctx, "")
	require.NoError(t, err)
	// Add http:// protocol for proper URL formatting
	apiEndpoint = "http://" + apiEndpoint
	t.Logf("API Gateway available at: %s", apiEndpoint)

	// Inspect logs from API Gateway to debug connection issues
	logs, err := apiGatewayC.Logs(ctx)
	if err == nil {
		logContent, _ := io.ReadAll(logs)
		t.Logf("API Gateway logs:\n%s", string(logContent))
	} else {
		t.Logf("Error retrieving API Gateway logs: %v", err)
	}
	defer logs.Close()

	// Wait for all services to be fully initialized and healthy
	if !utils.VerifyServicesInitialized(t, apiEndpoint) {
		t.Fatal("Services failed to initialize")
	}

	// Capture container logs before running tests to establish baseline
	utils.CaptureContainerLogs(ctx, t, "API Gateway Baseline Logs", apiGatewayC)
	utils.CaptureContainerLogs(ctx, t, "Auth Service Baseline Logs", authServiceC)
	utils.CaptureContainerLogs(ctx, t, "Odds Service Baseline Logs", oddsServiceC)

	// Default match ID to use for update and delete operations in individual tests
	defaultMatchID := "1" // Using a placeholder ID since we can't guarantee a specific ID exists

	// Run individual endpoint tests against the API Gateway
	t.Run("IndividualEndpointTests", func(t *testing.T) {
		t.Run("Login", func(t *testing.T) {
			t.Log("=== Starting Login Tests ===")
			// Capture current container state before test
			utils.CaptureContainerLogs(ctx, t, "Before Login Test - Auth Service", authServiceC)
			// Run the login test
			utils.RunLogin(t, apiEndpoint)
			// Capture container logs after test to see what happened
			utils.CaptureContainerLogs(ctx, t, "After Login Test - API Gateway", apiGatewayC)
			utils.CaptureContainerLogs(ctx, t, "After Login Test - Auth Service", authServiceC)
		})

		t.Run("CreateOdds", func(t *testing.T) {
			t.Log("=== Starting Create Odds Tests ===")
			token := utils.ObtainAuthToken(t, apiEndpoint)
			utils.TestCreateOddsEndpoint(t, apiEndpoint, token)
		})

		t.Run("ReadOdds", func(t *testing.T) {
			t.Log("=== Starting Read Odds Tests ===")
			token := utils.ObtainAuthToken(t, apiEndpoint)
			utils.TestReadOddsEndpoint(t, apiEndpoint, token)
		})

		t.Run("UpdateOdds", func(t *testing.T) {
			t.Log("=== Starting Update Odds Tests ===")
			token := utils.ObtainAuthToken(t, apiEndpoint)
			utils.TestUpdateOddsEndpoint(t, apiEndpoint, token, defaultMatchID)
		})

		t.Run("DeleteOdds", func(t *testing.T) {
			t.Log("=== Starting Delete Odds Tests ===")
			token := utils.ObtainAuthToken(t, apiEndpoint)
			utils.TestDeleteOddsEndpoint(t, apiEndpoint, token, defaultMatchID)
		})
	})

	// Run a complete user journey test that simulates a real user
	t.Run("E2EUserJourneyTest", func(t *testing.T) {
		utils.TestFullUserJourney(t, apiEndpoint)
	})
}



