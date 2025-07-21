package adapters

import (
	"context"
	"errors"
	"testing"

	"github.com/Businge931/sba-api-gateway/internal/api/models"
	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the AuthService interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RegisterResponse), args.Error(1)
}

func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.VerifyTokenResponse), args.Error(1)
}

func (m *MockAuthService) RequestPasswordReset(ctx context.Context, req *domain.RequestPasswordResetRequest) (*domain.RequestPasswordResetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RequestPasswordResetResponse), args.Error(1)
}

func (m *MockAuthService) ChangePassword(ctx context.Context, req *domain.ChangePasswordRequest) (*domain.ChangePasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChangePasswordResponse), args.Error(1)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, req *domain.ResetPasswordRequest) (*domain.ResetPasswordResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ResetPasswordResponse), args.Error(1)
}

func TestAuthServiceAdapter_Register(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		setupMock      func(*MockAuthService)
		req            *models.RegisterRequest
		expected       *models.RegisterResponse
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful registration",
			setupMock: func(m *MockAuthService) {
				m.On("Register", mock.Anything, &domain.RegisterRequest{
					Email:     "newuser@test.com",
					Password:  "SecurePass123!",
					FirstName: "newuser",
					LastName:  "newuser",
				}).Return(&domain.RegisterResponse{
					Success: true,
					Message: "User registered successfully",
				}, nil)
			},
			req: &models.RegisterRequest{
				Email:     "newuser@test.com",
				Password:  "SecurePass123!",
				FirstName: "newuser",
				LastName:  "newuser",
			},
			expected: &models.RegisterResponse{
				Success: true,
				Message: "User registered successfully",
			},
			expectedError: false,
		},
		{
			name: "username already exists",
			setupMock: func(m *MockAuthService) {
				m.On("Register", mock.Anything, &domain.RegisterRequest{
					Email:    "existinguser@test.com",
					Password: "password123",
				}).Return((*domain.RegisterResponse)(nil), errors.New("username already exists"))
			},
			req: &models.RegisterRequest{
				Email:    "existinguser@test.com",
				Password: "password123",
			},
			expected:       nil,
			expectedError:  true,
			expectedErrMsg: "username already exists",
		},
		{
			name: "weak password",
			setupMock: func(m *MockAuthService) {
				m.On("Register", mock.Anything, &domain.RegisterRequest{
					Email:    "newuser@test.com",
					Password: "123",
				}).Return((*domain.RegisterResponse)(nil), errors.New("password is too weak"))
			},
			req: &models.RegisterRequest{
				Email:    "newuser@test.com",
				Password: "123", // Weak password
			},
			expected:       nil,
			expectedError:  true,
			expectedErrMsg: "password is too weak",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockAuthService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Create adapter with mock service
			adapter := NewAuthServiceAdapter(mockSvc)

			// Call the method under test
			result, err := adapter.Register(context.Background(), tt.req)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestAuthServiceAdapter_Login(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		setupMock      func(*MockAuthService)
		req            *models.LoginRequest
		expected       *models.LoginResponse
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful login",
			setupMock: func(m *MockAuthService) {
				m.On("Login", mock.Anything, &domain.LoginRequest{
					Email:    "testuser@test.com",
					Password: "testpass",
				}).Return(&domain.LoginResponse{
					Success: true,
					Message: "Login successful",
					Token:   "test-token-123",
				}, nil)
			},
			req: &models.LoginRequest{
				Email:    "testuser@test.com",
				Password: "testpass",
			},
			expected: &models.LoginResponse{
				Success: true,
				Message: "Login successful",
				Token:   "test-token-123",
			},
			expectedError: false,
		},
		{
			name: "invalid credentials",
			setupMock: func(m *MockAuthService) {
				m.On("Login", mock.Anything, &domain.LoginRequest{
					Email:    "wronguser",
					Password: "wrongpass",
				}).Return((*domain.LoginResponse)(nil), errors.New("invalid credentials"))
			},
			req: &models.LoginRequest{
				Email:    "wronguser",
				Password: "wrongpass",
			},
			expected:       nil,
			expectedError:  true,
			expectedErrMsg: "invalid credentials",
		},
		{
			name: "empty username",
			setupMock: func(m *MockAuthService) {
				m.On("Login", mock.Anything, &domain.LoginRequest{
					Email:    "",
					Password: "testpass",
				}).Return((*domain.LoginResponse)(nil), errors.New("username is required"))
			},
			req: &models.LoginRequest{
				Email:    "",
				Password: "testpass",
			},
			expected:       nil,
			expectedError:  true,
			expectedErrMsg: "username is required",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockAuthService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Create adapter with mock service
			adapter := NewAuthServiceAdapter(mockSvc)

			// Call the method under test
			result, err := adapter.Login(context.Background(), tt.req)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestAuthServiceAdapter_VerifyToken(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		setupMock      func(*MockAuthService)
		token          string
		expected       *models.VerifyTokenResponse
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "valid token",
			setupMock: func(m *MockAuthService) {
				m.On("VerifyToken", mock.Anything, "valid-token-123").Return(&domain.VerifyTokenResponse{
					Success: true,
					Message: "Token is valid",
				}, nil)
			},
			token: "valid-token-123",
			expected: &models.VerifyTokenResponse{
				Success: true,
				Message: "Token is valid",
			},
			expectedError: false,
		},
		{
			name: "invalid token",
			setupMock: func(m *MockAuthService) {
				m.On("VerifyToken", mock.Anything, "invalid-token").Return(
					(*domain.VerifyTokenResponse)(nil),
					errors.New("invalid or expired token"),
				)
			},
			token:          "invalid-token",
			expected:       nil,
			expectedError:  true,
			expectedErrMsg: "invalid or expired token",
		},
		{
			name: "empty token",
			setupMock: func(m *MockAuthService) {
				m.On("VerifyToken", mock.Anything, "").Return(
					(*domain.VerifyTokenResponse)(nil),
					errors.New("token is required"),
				)
			},
			token:          "",
			expected:       nil,
			expectedError:  true,
			expectedErrMsg: "token is required",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockAuthService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Create adapter with mock service
			adapter := NewAuthServiceAdapter(mockSvc)

			// Call the method under test
			result, err := adapter.VerifyToken(context.Background(), tt.token)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}
