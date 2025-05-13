package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// Helper function to initialize the databases required for testing
func CreateAndInitDatabases(ctx context.Context, t *testing.T, postgresContainer testcontainers.Container) error {
	// Get postgres version to confirm container is working properly
	t.Log("Checking PostgreSQL container functionality...")
	versionCmd := []string{"psql", "-U", "postgres", "-c", "SELECT version();"}  
	exitCode, stdout, stderr := postgresContainer.Exec(ctx, versionCmd)
	if exitCode != 0 {
		t.Logf("Error checking PostgreSQL version: %s", stderr)
		return fmt.Errorf("postgres container not responding correctly: %s", stderr)
	}
	t.Logf("PostgreSQL version: %s", stdout)

	// Drop existing databases if they exist (clean slate)
	t.Log("Dropping existing databases if they exist...")
	dropDbsCmd := []string{
		"psql", "-U", "postgres", "-c", 
		"DROP DATABASE IF EXISTS sba_user_accounts; DROP DATABASE IF EXISTS sba_odds;",
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, dropDbsCmd)
	if exitCode != 0 {
		t.Logf("Warning while dropping databases: %s", stderr)
		// Continue despite errors here
	}
	t.Logf("Drop databases result: %s", stdout)

	// Create the auth service database with explicit encoding
	t.Log("Creating auth service database...")
	createAuthDbCmd := []string{
		"psql", "-U", "postgres", "-c", 
		"CREATE DATABASE sba_user_accounts WITH ENCODING='UTF8' LC_COLLATE='en_US.utf8' LC_CTYPE='en_US.utf8';",
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, createAuthDbCmd)
	if exitCode != 0 {
		t.Logf("Error creating auth service database: %s", stderr)
		return fmt.Errorf("failed to create auth service database: %s", stderr)
	} 
	t.Logf("Create auth DB result: %s", stdout)

	// Create the odds service database with explicit encoding
	t.Log("Creating odds service database...")
	createOddsDbCmd := []string{
		"psql", "-U", "postgres", "-c", 
		"CREATE DATABASE sba_odds WITH ENCODING='UTF8' LC_COLLATE='en_US.utf8' LC_CTYPE='en_US.utf8';",
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, createOddsDbCmd)
	if exitCode != 0 {
		t.Logf("Error creating odds service database: %s", stderr)
		return fmt.Errorf("failed to create odds service database: %s", stderr)
	}
	t.Logf("Create odds DB result: %s", stdout)

	// Grant all privileges to postgres user
	t.Log("Setting database permissions...")
	grantAuthDbCmd := []string{
		"psql", "-U", "postgres", "-c", 
		"GRANT ALL PRIVILEGES ON DATABASE sba_user_accounts TO postgres;",
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, grantAuthDbCmd)
	if exitCode != 0 {
		t.Logf("Warning while setting auth DB permissions: %s", stderr)
		// Continue despite errors here
	}

	grantOddsDbCmd := []string{
		"psql", "-U", "postgres", "-c", 
		"GRANT ALL PRIVILEGES ON DATABASE sba_odds TO postgres;",
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, grantOddsDbCmd)
	if exitCode != 0 {
		t.Logf("Warning while setting odds DB permissions: %s", stderr)
		// Continue despite errors here
	}

	// Create user accounts table in auth service database
	t.Log("Creating auth tables...")
	createAuthTableCmd := []string{
		"psql", "-U", "postgres", "-d", "sba_user_accounts", "-c",
		`CREATE TABLE IF NOT EXISTS user_accounts (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, createAuthTableCmd)
	if exitCode != 0 {
		t.Logf("Error creating auth service tables: %s", stderr)
		return fmt.Errorf("failed to create auth service tables: %s", stderr)
	}
	t.Logf("Create auth table result: %s", stdout)

	// Create a test user in the auth service database with BCrypt hashed password 'password123'
	t.Log("Creating test user in auth database...")
	createTestUserCmd := []string{
		"psql", "-U", "postgres", "-d", "sba_user_accounts", "-c",
		`INSERT INTO user_accounts (username, password, email, first_name, last_name) 
		 VALUES ('testuser', '$2a$10$wPxvwdL1H1yCViIR1HYPb.3DMvUmeDvMcq.DTT9/4nHzYHgLFPHLu', 'test@example.com', 'Test', 'User') 
		 ON CONFLICT (username) DO NOTHING;`,
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, createTestUserCmd)
	if exitCode != 0 {
		t.Logf("Error creating test user: %s", stderr)
		return fmt.Errorf("failed to create test user: %s", stderr)
	}
	t.Logf("Create test user result: %s", stdout)

	// Verify user was created successfully
	verifyUserCmd := []string{
		"psql", "-U", "postgres", "-d", "sba_user_accounts", "-c",
		"SELECT id, username, email FROM user_accounts;",
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, verifyUserCmd)
	if exitCode != 0 {
		t.Logf("Error verifying test user: %s", stderr)
	} else {
		t.Logf("Verify test user result: %s", stdout)
	}

	// Create odds table in odds service database
	t.Log("Creating odds tables...")
	createOddsTableCmd := []string{
		"psql", "-U", "postgres", "-d", "sba_odds", "-c",
		`CREATE TABLE IF NOT EXISTS odds (
			id SERIAL PRIMARY KEY,
			event_name VARCHAR(255) NOT NULL,
			home_team VARCHAR(255) NOT NULL,
			away_team VARCHAR(255) NOT NULL,
			home_odds DECIMAL(10,2) NOT NULL,
			away_odds DECIMAL(10,2) NOT NULL,
			draw_odds DECIMAL(10,2),
			event_date TIMESTAMP NOT NULL,
			sport VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, createOddsTableCmd)
	if exitCode != 0 {
		t.Logf("Error creating odds service tables: %s", stderr)
		return fmt.Errorf("failed to create odds service tables: %s", stderr)
	}
	t.Logf("Create odds table result: %s", stdout)

	// Add some sample odds data for testing
	t.Log("Creating sample odds data...")
	createOddsDataCmd := []string{
		"psql", "-U", "postgres", "-d", "sba_odds", "-c",
		`INSERT INTO odds (event_name, home_team, away_team, home_odds, away_odds, draw_odds, event_date, sport) VALUES
		('Premier League Match', 'Arsenal', 'Chelsea', 2.10, 3.40, 3.20, '2025-06-01 15:00:00', 'Football'),
		('La Liga Match', 'Barcelona', 'Real Madrid', 1.80, 3.75, 3.60, '2025-06-02 20:00:00', 'Football'),
		('NBA Finals', 'Lakers', 'Celtics', 1.95, 1.85, NULL, '2025-06-03 18:30:00', 'Basketball')
		ON CONFLICT DO NOTHING;
		`,
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, createOddsDataCmd)
	if exitCode != 0 {
		t.Logf("Error adding sample odds data: %s", stderr)
		// Continue despite errors here, not critical
	} else {
		t.Logf("Create sample odds data result: %s", stdout)
	}

	// Verify odds data was created successfully
	verifyOddsCmd := []string{
		"psql", "-U", "postgres", "-d", "sba_odds", "-c",
		"SELECT id, event_name, home_team, away_team FROM odds;",
	}
	exitCode, stdout, stderr = postgresContainer.Exec(ctx, verifyOddsCmd)
	if exitCode != 0 {
		t.Logf("Error verifying odds data: %s", stderr)
	} else {
		t.Logf("Verify odds data result: %s", stdout)
	}

	// Wait for database initialization to complete
	t.Log("Database initialization completed successfully")
	time.Sleep(2 * time.Second)
	return nil
}