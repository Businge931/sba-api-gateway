package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/proto"
)

type OddsServer struct {
	client proto.OddsServiceClient
}

func NewOddsServer(conn *grpc.ClientConn) *OddsServer {
	return &OddsServer{client: proto.NewOddsServiceClient(conn)}
}

func (c *OddsServer) CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error) {
	protoReq := &proto.CreateOddsRequest{
		League:          req.League,
		HomeTeam:        req.HomeTeam,
		AwayTeam:        req.AwayTeam,
		HomeTeamWinOdds: float32(req.HomeTeamWinOdds),
		AwayTeamWinOdds: float32(req.AwayTeamWinOdds),
		DrawOdds:        float32(req.DrawOdds),
		GameDate:        req.GameDate,
	}

	res, err := c.client.CreateOdds(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &domain.CreateOddsResponse{
		Success: res.GetSuccess(),
		Message: res.GetMessage(),
		Details: res.GetDetails(),
	}, nil
}

func (c *OddsServer) ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error) {
	protoReq := &proto.ReadOddsRequest{
		League: req.League,
		Date:   req.Date,
	}

	res, err := c.client.ReadOdds(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	// Map the gRPC response to the domain response
	odds := make([]domain.CreateOddsRequest, len(res.Odds))
	for i, protoOdds := range res.GetOdds() {
		odds[i] = domain.CreateOddsRequest{
			League:          protoOdds.GetLeague(),
			HomeTeam:        protoOdds.GetHomeTeam(),
			AwayTeam:        protoOdds.GetAwayTeam(),
			HomeTeamWinOdds: float64(protoOdds.GetHomeTeamWinOdds()),
			AwayTeamWinOdds: float64(protoOdds.GetAwayTeamWinOdds()),
			DrawOdds:        float64(protoOdds.GetDrawOdds()),
			GameDate:        protoOdds.GameDate,
		}
	}

	return &domain.ReadOddsResponse{
		Odds:    odds,
		Details: res.GetDetails(),
	}, nil
}
func (c *OddsServer) UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error) {
	protoReq := &proto.UpdateOddsRequest{
		League:          req.League,
		HomeTeam:        req.HomeTeam,
		AwayTeam:        req.AwayTeam,
		HomeTeamWinOdds: float32(req.HomeTeamWinOdds),
		AwayTeamWinOdds: float32(req.AwayTeamWinOdds),
		DrawOdds:        float32(req.DrawOdds),
		GameDate:        req.GameDate,
	}

	res, err := c.client.UpdateOdds(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &domain.UpdateOddsResponse{
		Success: res.GetSuccess(),
		Message: res.GetMessage(),
		Details: res.GetDetails(),
	}, nil
}

func (c *OddsServer) DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error) {
	protoReq := &proto.DeleteOddsRequest{
		League:   req.League,
		HomeTeam: req.HomeTeam,
		AwayTeam: req.AwayTeam,
		GameDate: req.GameDate,
	}

	res, err := c.client.DeleteOdds(ctx, protoReq)
	if err != nil {
		return nil, err
	}

	return &domain.DeleteOddsResponse{
		Success: res.GetSuccess(),
		Message: res.GetMessage(),
		Details: res.GetDetails(),
	}, nil
}
