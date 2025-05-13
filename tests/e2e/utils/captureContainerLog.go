package utils

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

// Helper function to capture and log container logs for debugging
func CaptureContainerLogs(ctx context.Context, t *testing.T, label string, container testcontainers.Container) {
	t.Logf("===== %s =====", label)

	// Skip getting container ID since the API may have changed
	t.Logf("Capturing logs for container: %T", container)

	// Get container logs
	logs, err := container.Logs(ctx)
	if err != nil {
		t.Logf("Error getting container logs: %v", err)
		return
	}

	// Read all logs
	logContent, err := io.ReadAll(logs)
	if err != nil {
		t.Logf("Error reading log content: %v", err)
		return
	}

	// Truncate logs if they're too large
	logStr := string(logContent)
	if len(logStr) > 5000 {
		// Show last 5000 characters which are most relevant
		logStr = "..." + logStr[len(logStr)-5000:]
	}

	t.Logf("Container Logs:\n%s\n", logStr)

	// Check for common error patterns
	if strings.Contains(logStr, "error") || strings.Contains(logStr, "Error") || 
	   strings.Contains(logStr, "ERROR") || strings.Contains(logStr, "failed") || 
	   strings.Contains(logStr, "Failed") {
		t.Log("DETECTED ERROR PATTERNS IN LOGS (see above)")
	}

	t.Logf("===== END %s =====", label)
}
