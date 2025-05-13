package utils

import (
	"net/http"
	"testing"
)

// Define a custom http.RoundTripper for detailed request/response logging
type LoggingRoundTripper struct {
	t        *testing.T
	delegate http.RoundTripper
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	lrt.t.Logf("HTTP Request: %s %s", req.Method, req.URL.String())
	lrt.t.Log("Request Headers:")
	for key, values := range req.Header {
		for _, value := range values {
			lrt.t.Logf("%s: %s", key, value)
		}
	}

	// Forward the request to the actual transport
	resp, err := lrt.delegate.RoundTrip(req)
	if err != nil {
		lrt.t.Logf("HTTP Transport error: %v", err)
		return nil, err
	}

	lrt.t.Logf("HTTP Response: %d %s", resp.StatusCode, resp.Status)
	return resp, nil
}
