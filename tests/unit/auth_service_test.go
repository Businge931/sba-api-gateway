package unit

import (
	"context"
	"testing"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthServer is a mock of the AuthServer interface
type MockAuthServer struct {
	mock.Mock
}

// Login mocks the Login method
func (m *MockAuthServer) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.LoginResponse), args.Error(1)
}

// Register mocks the Register method
func (m *MockAuthServer) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.RegisterResponse), args.Error(1)
}

// VerifyToken mocks the VerifyToken method
func (m *MockAuthServer) VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*domain.VerifyTokenResponse), args.Error(1)
}

func TestAuthService(t *testing.T) {
	// Create a mock auth server
	mockServer := new(MockAuthServer)

	// Import the service
	service := NewAuthServiceForTesting(mockServer)

	// Test Login
	t.Run("Login", func(t *testing.T) {
		tests := []struct {
			name            string
			Email           string
			password        string
			expectedSuccess bool
			expectedToken   string
			expectedMessage string
			expectedError   error
		}{
			{
				name:            "successful login",
				Email:           "testuser@test.com",
				password:        "testpass",
				expectedSuccess: true,
				expectedToken:   "test-token",
				expectedMessage: "Login successful",
				expectedError:   nil,
			},
			{
				name:            "invalid credentials",
				Email:           "wronguser@test.com",
				password:        "wrongpass",
				expectedSuccess: false,
				expectedToken:   "",
				expectedMessage: "Invalid credentials",
				expectedError:   nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Set up expectations
				loginReq := &domain.LoginRequest{
					Email:    tc.Email,
					Password: tc.password,
				}
				expectedResponse := &domain.LoginResponse{
					Success: tc.expectedSuccess,
					Token:   tc.expectedToken,
					Message: tc.expectedMessage,
				}

				mockServer.On("Login", mock.Anything, loginReq).Return(expectedResponse, tc.expectedError).Once()

				// Call the method
				response, err := service.Login(context.Background(), loginReq)

				// Assert expectations
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, expectedResponse, response)
				}
			})
		}

		mockServer.AssertExpectations(t)
	})

	// Test Register
	t.Run("Register", func(t *testing.T) {
		tests := []struct {
			name            string
			username        string
			password        string
			expectedSuccess bool
			expectedMessage string
			expectedError   error
		}{
			{
				name:            "successful registration",
				username:        "newuser",
				password:        "newpass",
				expectedSuccess: true,
				expectedMessage: "Registration successful",
				expectedError:   nil,
			},
			{
				name:            "user already exists",
				username:        "existinguser",
				password:        "somepass",
				expectedSuccess: false,
				expectedMessage: "User already exists",
				expectedError:   nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Set up expectations
				registerReq := &domain.RegisterRequest{
					Email:    tc.username,
					Password: tc.password,
				}
				expectedResponse := &domain.RegisterResponse{
					Success: tc.expectedSuccess,
					Message: tc.expectedMessage,
				}

				mockServer.On("Register", mock.Anything, registerReq).Return(expectedResponse, tc.expectedError).Once()

				// Call the method
				response, err := service.Register(context.Background(), registerReq)

				// Assert expectations
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, expectedResponse, response)
				}
			})
		}

		mockServer.AssertExpectations(t)
	})

	// Test VerifyToken
	t.Run("VerifyToken", func(t *testing.T) {
		tests := []struct {
			name            string
			token           string
			expectedSuccess bool
			expectedMessage string
			expectedError   error
		}{
			{
				name:            "valid token",
				token:           "valid-token",
				expectedSuccess: true,
				expectedMessage: "Token is valid",
				expectedError:   nil,
			},
			{
				name:            "invalid token",
				token:           "invalid-token",
				expectedSuccess: false,
				expectedMessage: "Invalid token",
				expectedError:   nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Set up expectations
				expectedResponse := &domain.VerifyTokenResponse{
					Success: tc.expectedSuccess,
					Message: tc.expectedMessage,
				}

				mockServer.On("VerifyToken", mock.Anything, tc.token).Return(expectedResponse, tc.expectedError).Once()

				// Call the method
				response, err := service.VerifyToken(context.Background(), tc.token)

				// Assert expectations
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, expectedResponse, response)
				}
			})
		}

		mockServer.AssertExpectations(t)
	})
}

// AuthService is a subset of the service interface we need for testing
type AuthService interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error)
	VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error)
}

// AuthServer is the interface we need to mock
type AuthServer interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error)
	VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error)
}

// authService implements AuthService
type authService struct {
	server AuthServer
}

// NewAuthServiceForTesting creates a new AuthService for testing
func NewAuthServiceForTesting(server AuthServer) AuthService {
	return &authService{
		server: server,
	}
}

// Login handles login requests
func (s *authService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	return s.server.Login(ctx, req)
}

// Register handles registration requests
func (s *authService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	return s.server.Register(ctx, req)
}

// VerifyToken verifies a token
func (s *authService) VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	return s.server.VerifyToken(ctx, token)
}
