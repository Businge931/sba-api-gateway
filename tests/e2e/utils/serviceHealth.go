package utils

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func VerifyServicesInitialized(t *testing.T, gatewayEndpoint string) bool {
	t.Log("Waiting for services to fully initialize...")
	return VerifyServiceHealth(gatewayEndpoint + "/health", t)
}

// Helper function to verify a service health endpoint
func VerifyServiceHealth(healthEndpoint string, t *testing.T) bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Try the health endpoint a few times
	for i := 0; i < 5; i++ {
		resp, err := client.Get(healthEndpoint)
		if err != nil {
			t.Logf("Health check attempt %d failed: %v", i+1, err)
			time.Sleep(1 * time.Second)
			continue
		}
		
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Logf("Health check status: %d, body: %s", resp.StatusCode, body)
		if resp.StatusCode == http.StatusOK {
			return true
		}

		time.Sleep(2 * time.Second)
	}

	return false
}