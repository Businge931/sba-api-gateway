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
	service AuthService
}

func NewAuthService(service AuthService) *AuthServiceImpl {
	return &AuthServiceImpl{service: service}
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	return s.service.Login(ctx, req)
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	return s.service.Register(ctx, req)
}

func (s *AuthServiceImpl) VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	return s.service.VerifyToken(ctx, token)
}
