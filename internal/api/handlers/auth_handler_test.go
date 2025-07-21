package handlers

import (
	"bytes"
	"context"
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

// MockAuthService is a mock implementation of the AuthService interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (*models.VerifyTokenResponse, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.VerifyTokenResponse), args.Error(1)
}

func (m *MockAuthService) RequestPasswordReset(ctx context.Context, req *models.RequestPasswordResetRequest) (*models.RequestPasswordResetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RequestPasswordResetResponse), args.Error(1)
}

func (m *MockAuthService) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) (*models.ChangePasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ChangePasswordResponse), args.Error(1)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ResetPasswordResponse), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockAuthService)
		requestBody   interface{}
		wantStatus    int
		checkResponse func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "successful registration",
			setupMock: func(mas *MockAuthService) {
				mas.On("Register", mock.Anything, &models.RegisterRequest{
					Email:    "newuser@test.com",
					Password: "newpass123",
				}).Return(&models.RegisterResponse{
					Success: true,
					Message: "Registration successful",
				}, nil)
			},
			requestBody: map[string]string{
				"email":    "newuser@test.com",
				"password": "newpass123",
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response["success"].(bool))
				assert.Equal(t, "Registration successful", response["message"])
			},
		},
		{
			name:        "invalid request body",
			setupMock:   func(mas *MockAuthService) {},
			requestBody: "invalid-json",
			wantStatus:  http.StatusForbidden,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Contains(t, rr.Body.String(), "Invalid request")
			},
		},
		{
			name:      "missing registration data",
			setupMock: func(mas *MockAuthService) {},
			requestBody: map[string]string{
				"email":    "",
				"password": "",
			},
			wantStatus: http.StatusForbidden,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Contains(t, rr.Body.String(), "Invalid request data")
			},
		},
		{
			name: "registration failure - username taken",
			setupMock: func(mas *MockAuthService) {
				mas.On("Register", mock.Anything, mock.Anything).Return(
					nil,
					status.Error(codes.AlreadyExists, "username already taken"),
				)
			},
			requestBody: map[string]string{
				"email":    "existinguser@test.com",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.False(t, response["success"].(bool))
				assert.Contains(t, response["error"].(string), "username already taken")
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

			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.Register(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}

			mockAuthSvc.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(*MockAuthService)
		requestBody   interface{}
		wantStatus    int
		wantResponse  map[string]interface{}
		checkResponse func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			name: "successful login",
			setupMock: func(mas *MockAuthService) {
				mas.On("Login", mock.Anything, &models.LoginRequest{
					Email:    "testuser@test.com",
					Password: "testpass",
				}).Return(&models.LoginResponse{
					Success: true,
					Message: "Login successful",
					Token:   "test-token",
				}, nil)
			},
			requestBody: map[string]string{
				"email":    "testuser@test.com",
				"password": "testpass",
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response["success"].(bool))
				assert.Equal(t, "Login successful", response["message"])
				assert.Equal(t, "test-token", response["token"])
			},
		},
		{
			name:        "invalid request body",
			setupMock:   func(mas *MockAuthService) {},
			requestBody: "invalid-json",
			wantStatus:  http.StatusForbidden,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Contains(t, rr.Body.String(), "Invalid request")
			},
		},
		{
			name:      "missing credentials",
			setupMock: func(mas *MockAuthService) {},
			requestBody: map[string]string{
				"email":    "",
				"password": "",
			},
			wantStatus: http.StatusForbidden,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Contains(t, rr.Body.String(), "Invalid request data")
			},
		},
		{
			name: "authentication failure",
			setupMock: func(mas *MockAuthService) {
				mas.On("Login", mock.Anything, mock.Anything).Return(
					nil,
					status.Error(codes.Unauthenticated, "invalid credentials"),
				)
			},
			requestBody: map[string]string{
				"email":    "wronguser@test.com",
				"password": "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, rr *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.False(t, response["success"].(bool))
				assert.Contains(t, response["error"].(string), "invalid credentials")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockAuthSvc := &MockAuthService{}
			if tt.setupMock != nil {
				tt.setupMock(mockAuthSvc)
			}

			handler := NewAuthHandler(mockAuthSvc)

			// Create request
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call handler
			handler.Login(rr, req)

			// Assertions
			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}

			// Verify mock expectations
			mockAuthSvc.AssertExpectations(t)
		})
	}
}
