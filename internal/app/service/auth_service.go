
package service

import (
	"context"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
)

// AuthService defines the interface for authentication-related operations.
type AuthService interface {
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error)
	VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error)
}

// AuthServiceImpl implements the AuthService interface.
type AuthServiceImpl struct {
	client domain.AuthClient
}

func NewAuthService(client domain.AuthClient) *AuthServiceImpl {
	return &AuthServiceImpl{client: client}
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	return s.client.Login(ctx, req)
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	return s.client.Register(ctx, req)
}

func (s *AuthServiceImpl) VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	return s.client.VerifyToken(ctx, token)
}
