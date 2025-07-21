package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/proto"
)

type AuthServer struct {
	client proto.AuthServiceClient
}

func NewAuthServer(conn *grpc.ClientConn) *AuthServer {
	return &AuthServer{client: proto.NewAuthServiceClient(conn)}
}

func (c *AuthServer) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	protoReq := &proto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	res, err := c.client.Login(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.LoginResponse{
		Success: true,
		Token:   res.GetToken(),
		Message: "Login successful",
	}, nil
}

func (c *AuthServer) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	protoReq := &proto.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	res, err := c.client.Register(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.RegisterResponse{
		Success: true,
		Message: res.GetMessage(),
	}, nil
}

func (c *AuthServer) VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	res, err := c.client.VerifyToken(ctx, &proto.VerifyTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	return &domain.VerifyTokenResponse{
		Success: res.GetSuccess(),
		Message: res.GetMessage(),
	}, nil
}

func (c *AuthServer) RequestPasswordReset(ctx context.Context, req *domain.RequestPasswordResetRequest) (*domain.RequestPasswordResetResponse, error) {
	protoReq := &proto.RequestPasswordResetRequest{
		Email: req.Email,
	}
	res, err := c.client.RequestPasswordReset(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.RequestPasswordResetResponse{
		Success: res.GetSuccess(),
		Message: res.GetMessage(),
	}, nil
}

func (c *AuthServer) ChangePassword(ctx context.Context, req *domain.ChangePasswordRequest) (*domain.ChangePasswordResponse, error) {
	protoReq := &proto.ChangePasswordRequest{
		UserId:      req.UserID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	res, err := c.client.ChangePassword(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.ChangePasswordResponse{
		Success: res.GetSuccess(),
		Message: res.GetMessage(),
	}, nil
}

func (c *AuthServer) ResetPassword(ctx context.Context, req *domain.ResetPasswordRequest) (*domain.ResetPasswordResponse, error) {
	protoReq := &proto.ResetPasswordRequest{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}
	res, err := c.client.ResetPassword(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.ResetPasswordResponse{
		Success: res.GetSuccess(),
		Message: res.GetMessage(),
		UserID:  res.GetUserId(),
	}, nil
}
