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
		Username: req.Username,
		Password: req.Password,
	}
	res, err := c.client.Login(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.LoginResponse{
		Success: true, // Set success to true when login is successful
		Token: res.GetToken(),
		Message: "Login successful", // Add a success message
	}, nil
}

func (c *AuthServer) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	protoReq := &proto.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	}
	res, err := c.client.Register(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return &domain.RegisterResponse{
		Success: true, // Set success to true when registration is successful
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
