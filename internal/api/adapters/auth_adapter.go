package adapters

import (
	"context"

	"github.com/Businge931/sba-api-gateway/internal/api/models"
	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
)

// AuthServiceAdapter adapts the service.AuthService to handlers.AuthService
type AuthServiceAdapter struct {
	service service.AuthService
}

// NewAuthServiceAdapter creates a new adapter for the auth service
func NewAuthServiceAdapter(service service.AuthService) *AuthServiceAdapter {
	return &AuthServiceAdapter{
		service: service,
	}
}

// Login adapts between models.LoginRequest and domain.LoginRequest
func (a *AuthServiceAdapter) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Convert models.LoginRequest to domain.LoginRequest
	domainReq := &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call the service with the domain request
	domainRes, err := a.service.Login(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.LoginResponse to models.LoginResponse
	return &models.LoginResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
		Token:   domainRes.Token,
	}, nil
}

// Register adapts between models.RegisterRequest and domain.RegisterRequest
func (a *AuthServiceAdapter) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	// Convert models.RegisterRequest to domain.RegisterRequest
	domainReq := &domain.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Call the service with the domain request
	domainRes, err := a.service.Register(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.RegisterResponse to models.RegisterResponse
	return &models.RegisterResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
	}, nil
}

// VerifyToken adapts the service.AuthService.VerifyToken to handlers.AuthService.VerifyToken
func (a *AuthServiceAdapter) VerifyToken(ctx context.Context, token string) (*models.VerifyTokenResponse, error) {
	// Call the service with the token
	domainRes, err := a.service.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Convert domain.VerifyTokenResponse to models.VerifyTokenResponse
	return &models.VerifyTokenResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
	}, nil
}

// RequestPasswordReset adapts between models.RequestPasswordResetRequest and domain.RequestPasswordResetRequest
func (a *AuthServiceAdapter) RequestPasswordReset(ctx context.Context, req *models.RequestPasswordResetRequest) (*models.RequestPasswordResetResponse, error) {
	// Convert models.RequestPasswordResetRequest to domain.RequestPasswordResetRequest
	domainReq := &domain.RequestPasswordResetRequest{
		Email: req.Email,
	}

	// Call the service with the domain request
	domainRes, err := a.service.RequestPasswordReset(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.RequestPasswordResetResponse to models.RequestPasswordResetResponse
	return &models.RequestPasswordResetResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
	}, nil
}

// ChangePassword adapts between models.ChangePasswordRequest and domain.ChangePasswordRequest
func (a *AuthServiceAdapter) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) (*models.ChangePasswordResponse, error) {
	// Convert models.ChangePasswordRequest to domain.ChangePasswordRequest
	domainReq := &domain.ChangePasswordRequest{
		UserID:      req.UserID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	// Call the service with the domain request
	domainRes, err := a.service.ChangePassword(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.ChangePasswordResponse to models.ChangePasswordResponse
	return &models.ChangePasswordResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
	}, nil
}

// ResetPassword adapts between models.ResetPasswordRequest and domain.ResetPasswordRequest
func (a *AuthServiceAdapter) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error) {
	// Convert models.ResetPasswordRequest to domain.ResetPasswordRequest
	domainReq := &domain.ResetPasswordRequest{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}

	// Call the service with the domain request
	domainRes, err := a.service.ResetPassword(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.ResetPasswordResponse to models.ResetPasswordResponse
	return &models.ResetPasswordResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
		UserID:  domainRes.UserID,
	}, nil
}
