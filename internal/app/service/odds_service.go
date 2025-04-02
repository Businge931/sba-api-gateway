package service

import (
	"context"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
)

type OddsService interface {
	CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error)
	ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error)
	UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error)
	DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error)
}

type oddsService struct {
	client domain.OddsClient
}

func NewOddsService(client domain.OddsClient) *oddsService {
	return &oddsService{client: client}
}

func (s *oddsService) CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error) {
	return s.client.CreateOdds(ctx, req)
}

func (s *oddsService) ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error) {
	return s.client.ReadOdds(ctx, req)
}

func (s *oddsService) UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error) {
	return s.client.UpdateOdds(ctx, req)
}

func (s *oddsService) DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error) {
	return s.client.DeleteOdds(ctx, req)
}