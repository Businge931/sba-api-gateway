package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Businge931/sba-api-gateway/internal/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestVerifyTokenMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupMock     func(*MockAuthService)
		token         string
		handler       http.HandlerFunc
		wantStatus    int
		checkResponse func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "valid token",
			setupMock: func(mas *MockAuthService) {
				mas.On("VerifyToken", mock.Anything, "valid-token").Return(
					&models.VerifyTokenResponse{
						Success: true,
						Message: "Token is valid",
					},
					nil,
				)
			},
			token: "valid-token",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"success"}`))
			}),
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, `{"status":"success"}`, rr.Body.String())
			},
		},
		{
			name:       "missing token",
			setupMock: func(mas *MockAuthService) {},
			token:     "",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Contains(t, rr.Body.String(), "Missing token")
			},
		},
		{
			name: "invalid token",
			setupMock: func(mas *MockAuthService) {
				mas.On("VerifyToken", mock.Anything, "invalid-token").Return(
					&models.VerifyTokenResponse{
						Success: false,
						Message: "invalid token",
					},
					nil,
				)
			},
			token: "invalid-token",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Contains(t, rr.Body.String(), "invalid token")
			},
		},
		{
			name: "token verification error",
			setupMock: func(mas *MockAuthService) {
				mas.On("VerifyToken", mock.Anything, "error-token").Return(
					nil,
					status.Error(codes.Internal, "internal server error"),
				)
			},
			token: "error-token",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			wantStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.False(t, response["success"].(bool))
				assert.Contains(t, response["error"].(string), "internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthSvc := &MockAuthService{}
			if tt.setupMock != nil {
				tt.setupMock(mockAuthSvc)
			}

			handler := NewAuthHandler(mockAuthSvc)

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			rr := httptest.NewRecorder()

			// Create a test handler to wrap with the middleware
			middleware := handler.VerifyTokenMiddleware(tt.handler)
			middleware.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}

			mockAuthSvc.AssertExpectations(t)
		})
	}
}
