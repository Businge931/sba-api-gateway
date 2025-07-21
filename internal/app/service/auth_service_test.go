package service

import (
	"context"
	"errors"
	"testing"

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

func TestAuthServiceImpl_Login(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		mockSetup     func(*MockAuthService)
		args          *domain.LoginRequest
		want          *domain.LoginResponse
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful login",
			mockSetup: func(mas *MockAuthService) {
				mas.On("Login", mock.Anything, &domain.LoginRequest{
					Email:    "testuser@test.com",
					Password: "testpass",
				}).Return(&domain.LoginResponse{
					Success: true,
					Token:   "test-token",
					Message: "Login successful",
				}, nil)
			},
			args: &domain.LoginRequest{
				Email:    "testuser@test.com",
				Password: "testpass",
			},
			want: &domain.LoginResponse{
				Success: true,
				Token:   "test-token",
				Message: "Login successful",
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			mockSetup: func(mas *MockAuthService) {
				mas.On("Login", mock.Anything, &domain.LoginRequest{
					Email:    "wronguser@test.com",
					Password: "wrongpass",
				}).Return((*domain.LoginResponse)(nil), errors.New("invalid credentials"))
			},
			args: &domain.LoginRequest{
				Email:    "wronguser@test.com",
				Password: "wrongpass",
			},
			want:          nil,
			wantErr:       true,
			expectedError: errors.New("invalid credentials"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthSvc := new(MockAuthService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockAuthSvc)
			}

			svc := &AuthServiceImpl{
				service: mockAuthSvc,
			}

			got, err := svc.Login(context.Background(), tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockAuthSvc.AssertExpectations(t)
		})
	}
}

func TestAuthServiceImpl_Register(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		mockSetup     func(*MockAuthService)
		args          *domain.RegisterRequest
		want          *domain.RegisterResponse
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful registration",
			mockSetup: func(mas *MockAuthService) {
				mas.On("Register", mock.Anything, &domain.RegisterRequest{
					Email:    "newuser@test.com",
					Password: "newpass",
				}).Return(&domain.RegisterResponse{
					Success: true,
					Message: "Registration successful",
				}, nil)
			},
			args: &domain.RegisterRequest{
				Email:    "newuser@test.com",
				Password: "newpass",
			},
			want: &domain.RegisterResponse{
				Success: true,
				Message: "Registration successful",
			},
			wantErr: false,
		},
		{
			name: "duplicate username",
			mockSetup: func(mas *MockAuthService) {
				mas.On("Register", mock.Anything, &domain.RegisterRequest{
					Email:    "existinguser@test.com",
					Password: "password",
				}).Return((*domain.RegisterResponse)(nil), errors.New("username already exists"))
			},
			args: &domain.RegisterRequest{
				Email:    "existinguser@test.com",
				Password: "password",
			},
			want:          nil,
			wantErr:       true,
			expectedError: errors.New("username already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthSvc := new(MockAuthService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockAuthSvc)
			}

			svc := &AuthServiceImpl{
				service: mockAuthSvc,
			}

			got, err := svc.Register(context.Background(), tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockAuthSvc.AssertExpectations(t)
		})
	}
}

func TestAuthServiceImpl_VerifyToken(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		mockSetup     func(*MockAuthService)
		token         string
		want          *domain.VerifyTokenResponse
		wantErr       bool
		expectedError error
	}{
		{
			name: "valid token",
			mockSetup: func(mas *MockAuthService) {
				mas.On("VerifyToken", mock.Anything, "valid-token").Return(&domain.VerifyTokenResponse{
					Success: true,
					Message: "Token is valid",
				}, nil)
			},
			token: "valid-token",
			want: &domain.VerifyTokenResponse{
				Success: true,
				Message: "Token is valid",
			},
			wantErr: false,
		},
		{
			name: "invalid token",
			mockSetup: func(mas *MockAuthService) {
				mas.On("VerifyToken", mock.Anything, "invalid-token").Return((*domain.VerifyTokenResponse)(nil), errors.New("invalid token"))
			},
			token:         "invalid-token",
			want:          nil,
			wantErr:       true,
			expectedError: errors.New("invalid token"),
		},
		{
			name: "expired token",
			mockSetup: func(mas *MockAuthService) {
				mas.On("VerifyToken", mock.Anything, "expired-token").Return((*domain.VerifyTokenResponse)(nil), errors.New("token has expired"))
			},
			token:         "expired-token",
			want:          nil,
			wantErr:       true,
			expectedError: errors.New("token has expired"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthSvc := new(MockAuthService)
			if tt.mockSetup != nil {
				tt.mockSetup(mockAuthSvc)
			}

			svc := &AuthServiceImpl{
				service: mockAuthSvc,
			}

			got, err := svc.VerifyToken(context.Background(), tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.EqualError(t, err, tt.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockAuthSvc.AssertExpectations(t)
		})
	}
}
