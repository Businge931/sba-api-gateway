package domain

import (
	"context"
)

type OddsClient interface {
	CreateOdds(ctx context.Context, req *CreateOddsRequest) (*CreateOddsResponse, error)
	ReadOdds(ctx context.Context, req *ReadOddsRequest) (*ReadOddsResponse, error)
	UpdateOdds(ctx context.Context, req *UpdateOddsRequest) (*UpdateOddsResponse, error)
	DeleteOdds(ctx context.Context, req *DeleteOddsRequest) (*DeleteOddsResponse, error)
}

type CreateOddsRequest struct {
	League          string
	HomeTeam        string
	AwayTeam        string
	HomeTeamWinOdds float64
	AwayTeamWinOdds float64
	DrawOdds        float64
	GameDate        string
}

type CreateOddsResponse struct {
	Success bool
	Message string
	Details string
}

type ReadOddsRequest struct {
	League string
	Date   string
}

type ReadOddsResponse struct {
	Odds    []CreateOddsRequest
	Details string
}

type UpdateOddsRequest struct {
	League          string
	HomeTeam        string
	AwayTeam        string
	HomeTeamWinOdds float64
	AwayTeamWinOdds float64
	DrawOdds        float64
	GameDate        string
}

type UpdateOddsResponse struct {
	Success bool
	Message string
	Details string
}

type DeleteOddsRequest struct {
	League   string
	HomeTeam string
	AwayTeam string
	GameDate string
}

type DeleteOddsResponse struct {
	Success bool
	Message string
	Details string
}
