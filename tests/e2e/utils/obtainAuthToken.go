package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Helper function to get an auth token for testing protected endpoints
func ObtainAuthToken(t *testing.T, apiEndpoint string) string {
	// Create login request with our test user that was inserted during database setup
	loginReq := map[string]string{
		"username": "test@example.com", 
		"password": "password123",
	}
	payload, err := json.Marshal(loginReq)
	require.NoError(t, err)

    // Send login request
    t.Log("Obtaining authentication token...")
    resp, err := http.Post(fmt.Sprintf("%s/login", apiEndpoint), "application/json", bytes.NewBuffer(payload))
    if err != nil {
        t.Fatalf("Error making auth request: %v", err)
    }
    defer resp.Body.Close()

    // Read response body for better error handling
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        t.Fatalf("Error reading auth response: %v", err)
    }

    // Check response status
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Failed to get auth token. Status: %d, Body: %s", resp.StatusCode, string(respBody))
    }

    // Parse the login response
    var loginResp map[string]interface{}
    err = json.Unmarshal(respBody, &loginResp)
    if err != nil {
        t.Fatalf("Error parsing auth response: %v", err)
    }

    // Extract and return the token
    token, ok := loginResp["token"].(string)
    if !ok || token == "" {
        t.Fatalf("Invalid or missing token in response: %v", loginResp)
    }

    t.Logf("Successfully obtained auth token (length: %d)", len(token))
    return token
}
