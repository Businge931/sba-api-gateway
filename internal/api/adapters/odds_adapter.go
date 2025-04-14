package adapters

import (
	"context"

	"github.com/Businge931/sba-api-gateway/internal/api/models"
	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
)

// OddsServiceAdapter adapts the service.OddsService to the handlers.OddsHandlerService interface
type OddsServiceAdapter struct {
	service service.OddsService
}

// NewOddsServiceAdapter creates a new adapter for the odds service
func NewOddsServiceAdapter(service service.OddsService) *OddsServiceAdapter {
	return &OddsServiceAdapter{
		service: service,
	}
}

// CreateOdds adapts between models.CreateOddsRequest and domain.CreateOddsRequest
func (a *OddsServiceAdapter) CreateOdds(ctx context.Context, req *models.CreateOddsRequest) (*models.CreateOddsResponse, error) {
	// Convert models.CreateOddsRequest to domain.CreateOddsRequest
	domainReq := &domain.CreateOddsRequest{
		League:          req.League,
		GameDate:        req.GameDate,
		HomeTeam:        req.HomeTeam,
		AwayTeam:        req.AwayTeam,
		HomeTeamWinOdds: req.HomeTeamWinOdds,
		AwayTeamWinOdds: req.AwayTeamWinOdds,
		DrawOdds:        req.DrawOdds,
	}

	// Call the service with the domain request
	domainRes, err := a.service.CreateOdds(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.CreateOddsResponse to models.CreateOddsResponse
	return &models.CreateOddsResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
	}, nil
}

// ReadOdds adapts between models.ReadOddsRequest and domain.ReadOddsRequest
func (a *OddsServiceAdapter) ReadOdds(ctx context.Context, req *models.ReadOddsRequest) (*models.ReadOddsResponse, error) {
	// Convert models.ReadOddsRequest to domain.ReadOddsRequest
	domainReq := &domain.ReadOddsRequest{
		League: req.League,
		Date:   req.Date,
	}

	// Call the service with the domain request
	domainRes, err := a.service.ReadOdds(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.ReadOddsResponse to models.ReadOddsResponse
	// Map the domain odds to models odds (we need to convert each item)
	modelOdds := make([]models.CreateOddsRequest, len(domainRes.Odds))
	for i, odds := range domainRes.Odds {
		modelOdds[i] = models.CreateOddsRequest{
			League:          odds.League,
			GameDate:        odds.GameDate,
			HomeTeam:        odds.HomeTeam,
			AwayTeam:        odds.AwayTeam,
			HomeTeamWinOdds: odds.HomeTeamWinOdds,
			AwayTeamWinOdds: odds.AwayTeamWinOdds,
			DrawOdds:        odds.DrawOdds,
		}
	}

	return &models.ReadOddsResponse{
		Odds:    modelOdds,
		Details: domainRes.Details,
	}, nil
}

// UpdateOdds adapts between models.UpdateOddsRequest and domain.UpdateOddsRequest
func (a *OddsServiceAdapter) UpdateOdds(ctx context.Context, req *models.UpdateOddsRequest) (*models.UpdateOddsResponse, error) {
	// Convert models.UpdateOddsRequest to domain.UpdateOddsRequest
	domainReq := &domain.UpdateOddsRequest{
		League:          req.League,
		GameDate:        req.GameDate,
		HomeTeam:        req.HomeTeam,
		AwayTeam:        req.AwayTeam,
		HomeTeamWinOdds: req.HomeTeamWinOdds,
		AwayTeamWinOdds: req.AwayTeamWinOdds,
		DrawOdds:        req.DrawOdds,
	}

	// Call the service with the domain request
	domainRes, err := a.service.UpdateOdds(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.UpdateOddsResponse to models.UpdateOddsResponse
	return &models.UpdateOddsResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
	}, nil
}

// DeleteOdds adapts between models.DeleteOddsRequest and domain.DeleteOddsRequest
func (a *OddsServiceAdapter) DeleteOdds(ctx context.Context, req *models.DeleteOddsRequest) (*models.DeleteOddsResponse, error) {
	// Convert models.DeleteOddsRequest to domain.DeleteOddsRequest
	domainReq := &domain.DeleteOddsRequest{
		League:   req.League,
		GameDate: req.GameDate,
	}

	// Call the service with the domain request
	domainRes, err := a.service.DeleteOdds(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain.DeleteOddsResponse to models.DeleteOddsResponse
	return &models.DeleteOddsResponse{
		Success: domainRes.Success,
		Message: domainRes.Message,
	}, nil
}
