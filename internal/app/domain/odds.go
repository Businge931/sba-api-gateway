package domain

import (
	"context"
)

type (
	OddsClient interface {
		CreateOdds(ctx context.Context, req *CreateOddsRequest) (*CreateOddsResponse, error)
		ReadOdds(ctx context.Context, req *ReadOddsRequest) (*ReadOddsResponse, error)
		UpdateOdds(ctx context.Context, req *UpdateOddsRequest) (*UpdateOddsResponse, error)
		DeleteOdds(ctx context.Context, req *DeleteOddsRequest) (*DeleteOddsResponse, error)
	}

	CreateOddsRequest struct {
		League          string
		HomeTeam        string
		AwayTeam        string
		HomeTeamWinOdds float64
		AwayTeamWinOdds float64
		DrawOdds        float64
		GameDate        string
	}

	CreateOddsResponse struct {
		Success bool
		Message string
		Details string
	}

	ReadOddsRequest struct {
		League string
		Date   string
	}

	ReadOddsResponse struct {
		Odds    []CreateOddsRequest
		Details string
	}

	UpdateOddsRequest struct {
		League          string
		HomeTeam        string
		AwayTeam        string
		HomeTeamWinOdds float64
		AwayTeamWinOdds float64
		DrawOdds        float64
		GameDate        string
	}

	UpdateOddsResponse struct {
		Success bool
		Message string
		Details string
	}

	DeleteOddsRequest struct {
		League   string
		HomeTeam string
		AwayTeam string
		GameDate string
	}

	DeleteOddsResponse struct {
		Success bool
		Message string
		Details string
	}
)
