package utils

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

// Helper function to verify that the required databases exist
func VerifyDatabasesExist(ctx context.Context, t *testing.T, postgresContainer testcontainers.Container) error {
	// Check auth service database
	authCheckCmd := []string{
		"psql", "-U", "postgres", "-c", "SELECT datname FROM pg_database WHERE datname='sba_user_accounts';",
	}
	
	// Execute command to check if auth database exists
	authExitCode, authOut, err := postgresContainer.Exec(ctx, authCheckCmd)
	if err != nil {
		t.Logf("Error executing auth database check: %v", err)
		return fmt.Errorf("error executing auth database check: %w", err)
	}
	
	if authExitCode != 0 {
		t.Logf("Auth service database check failed with exit code: %d", authExitCode)
		return fmt.Errorf("auth service database check failed with exit code: %d", authExitCode)
	}
	
	// Read output to check if database exists
	authResult, err := io.ReadAll(authOut)
	if err != nil {
		t.Logf("Error reading auth check output: %v", err)
		return fmt.Errorf("error reading auth check output: %w", err)
	}
	
	// Check if database exists in output
	if !strings.Contains(string(authResult), "sba_user_accounts") {
		t.Logf("Auth service database does not exist")
		return fmt.Errorf("auth service database does not exist")
	}
	
	t.Log("Auth service database verified")

	// Check odds service database
	oddsCheckCmd := []string{
		"psql", "-U", "postgres", "-c", "SELECT datname FROM pg_database WHERE datname='sba_odds';",
	}
	
	// Execute command to check if odds database exists
	oddsExitCode, oddsOut, err := postgresContainer.Exec(ctx, oddsCheckCmd)
	if err != nil {
		t.Logf("Error executing odds database check: %v", err)
		return fmt.Errorf("error executing odds database check: %w", err)
	}
	
	if oddsExitCode != 0 {
		t.Logf("Odds service database check failed with exit code: %d", oddsExitCode)
		return fmt.Errorf("odds service database check failed with exit code: %d", oddsExitCode)
	}
	
	// Read output to check if database exists
	oddsResult, err := io.ReadAll(oddsOut)
	if err != nil {
		t.Logf("Error reading odds check output: %v", err)
		return fmt.Errorf("error reading odds check output: %w", err)
	}
	
	// Check if database exists in output
	if !strings.Contains(string(oddsResult), "sba_odds") {
		t.Logf("Odds service database does not exist")
		return fmt.Errorf("odds service database does not exist")
	}
	
	t.Log("Odds service database verified")
	return nil
}