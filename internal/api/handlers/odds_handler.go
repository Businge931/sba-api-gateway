package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Businge931/sba-api-gateway/internal/api/models"
	"github.com/Businge931/sba-api-gateway/proto"
)

// OddsHandlerService defines the interface for odds operations at the handler level
type OddsHandlerService interface {
	CreateOdds(ctx context.Context, req *models.CreateOddsRequest) (*models.CreateOddsResponse, error)
	ReadOdds(ctx context.Context, req *models.ReadOddsRequest) (*models.ReadOddsResponse, error)
	UpdateOdds(ctx context.Context, req *models.UpdateOddsRequest) (*models.UpdateOddsResponse, error)
	DeleteOdds(ctx context.Context, req *models.DeleteOddsRequest) (*models.DeleteOddsResponse, error)
}

type OddsHandler struct {
	oddsService OddsHandlerService
}

func NewOddsHandler(oddsService OddsHandlerService) *OddsHandler {
	return &OddsHandler{oddsService: oddsService}
}

func validateLeagueAndDate(league, date string) (string, bool) {
	// Convert league name to lowercase for case-insensitive check
	if strings.ToLower(league) != "english premier league" {
		return `{"details": "League must be 'english premier league'"}`, false
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return `{"details": "Invalid date format. Use YYYY-MM-DD"}`, false
	}

	return "", true
}

// handleValidationError writes a validation error response
func handleValidationError(w http.ResponseWriter, errMsg string) {
	http.Error(w, errMsg, http.StatusForbidden)
}

// handleOddsRequest is a helper function to handle common logic for odds requests
func (h *OddsHandler) handleOddsRequest(
	w http.ResponseWriter,
	r *http.Request,
	serviceCall func(context.Context, interface{}) (interface{}, error),
	req interface{},
) {
	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handleValidationError(w, `{"details": "Invalid request"}`)
		return
	}

	// Validate league and date
	var league, date string
	switch r := req.(type) {
	case *proto.CreateOddsRequest:
		league, date = r.GetLeague(), r.GetGameDate()
	case *proto.ReadOddsRequest:
		league, date = r.GetLeague(), r.GetDate()
	case *proto.UpdateOddsRequest:
		league, date = r.GetLeague(), r.GetGameDate()
	case *proto.DeleteOddsRequest:
		league, date = r.GetLeague(), r.GetGameDate()
	}

	if errMsg, valid := validateLeagueAndDate(league, date); !valid {
		handleValidationError(w, errMsg)
		return
	}

	// Convert simple date format to RFC3339 format expected by CRUD service
	if parsedDate, err := time.Parse("2006-01-02", date); err == nil {
		// Set time to noon UTC to ensure it's a valid future time
		rfc3339Date := parsedDate.Format(time.RFC3339)

		// Update the request with the formatted date
		switch r := req.(type) {
		case *proto.CreateOddsRequest:
			r.GameDate = rfc3339Date
			// Convert league name to proper title case for CRUD service
			r.League = "English Premier League"
		case *proto.ReadOddsRequest:
			r.Date = rfc3339Date
			r.League = "English Premier League"
		case *proto.UpdateOddsRequest:
			r.GameDate = rfc3339Date
			r.League = "English Premier League"
		case *proto.DeleteOddsRequest:
			r.GameDate = rfc3339Date
			r.League = "English Premier League"
		}
	}

	// Call the service method
	res, err := serviceCall(r.Context(), req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

// CreateOdds handles the creation of odds
func (h *OddsHandler) CreateOdds(w http.ResponseWriter, r *http.Request) {
	wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
		protoReq := req.(*proto.CreateOddsRequest)
		domainReq := &models.CreateOddsRequest{
			League:          protoReq.League,
			GameDate:        protoReq.GameDate,
			HomeTeam:        protoReq.HomeTeam,
			AwayTeam:        protoReq.AwayTeam,
			HomeTeamWinOdds: float64(protoReq.HomeTeamWinOdds),
			AwayTeamWinOdds: float64(protoReq.AwayTeamWinOdds),
			DrawOdds:        float64(protoReq.DrawOdds),
		}
		return h.oddsService.CreateOdds(ctx, domainReq)
	}
	h.handleOddsRequest(w, r, wrapper, &proto.CreateOddsRequest{})
}

// ReadOdds handles the reading of odds
func (h *OddsHandler) ReadOdds(w http.ResponseWriter, r *http.Request) {
	// For GET requests, parse query parameters instead of request body
	if r.Method == http.MethodGet {
		league := r.URL.Query().Get("league")
		date := r.URL.Query().Get("date")

		// Validate league and date
		if errMsg, valid := validateLeagueAndDate(league, date); !valid {
			handleValidationError(w, errMsg)
			return
		}

		// Convert simple date format to RFC3339 format expected by CRUD service
		rfc3339Date := date
		if parsedDate, err := time.Parse("2006-01-02", date); err == nil {
			rfc3339Date = parsedDate.Format(time.RFC3339)
		}

		// Create the request
		protoReq := &proto.ReadOddsRequest{
			League: "English Premier League", // Proper title case
			Date:   rfc3339Date,
		}

		// Call the service
		domainReq := &models.ReadOddsRequest{
			League: protoReq.League,
			Date:   protoReq.Date,
		}

		res, err := h.oddsService.ReadOdds(r.Context(), domainReq)
		if err != nil {
			handleGRPCError(w, err)
			return
		}

		// Write the response
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// For POST requests or other methods, use the existing handleOddsRequest function
		wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
			protoReq := req.(*proto.ReadOddsRequest)
			domainReq := &models.ReadOddsRequest{
				League: protoReq.League,
				Date:   protoReq.Date,
			}
			return h.oddsService.ReadOdds(ctx, domainReq)
		}
		h.handleOddsRequest(w, r, wrapper, &proto.ReadOddsRequest{})
	}
}

// UpdateOdds handles the updating of odds
func (h *OddsHandler) UpdateOdds(w http.ResponseWriter, r *http.Request) {
	wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
		protoReq := req.(*proto.UpdateOddsRequest)
		domainReq := &models.UpdateOddsRequest{
			League:          protoReq.League,
			GameDate:        protoReq.GameDate,
			HomeTeam:        protoReq.HomeTeam,
			AwayTeam:        protoReq.AwayTeam,
			HomeTeamWinOdds: float64(protoReq.HomeTeamWinOdds),
			AwayTeamWinOdds: float64(protoReq.AwayTeamWinOdds),
			DrawOdds:        float64(protoReq.DrawOdds),
		}
		return h.oddsService.UpdateOdds(ctx, domainReq)
	}
	h.handleOddsRequest(w, r, wrapper, &proto.UpdateOddsRequest{})
}

// DeleteOdds handles the deletion of odds
func (h *OddsHandler) DeleteOdds(w http.ResponseWriter, r *http.Request) {
	wrapper := func(ctx context.Context, req interface{}) (interface{}, error) {
		protoReq := req.(*proto.DeleteOddsRequest)
		domainReq := &models.DeleteOddsRequest{
			League:   protoReq.League,
			GameDate: protoReq.GameDate,
		}
		return h.oddsService.DeleteOdds(ctx, domainReq)
	}
	h.handleOddsRequest(w, r, wrapper, &proto.DeleteOddsRequest{})
}
