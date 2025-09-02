package handlers

import (
	"context"
	"encoding/json"

	// "io"
	"net/http"
	"strings"
	"time"

	"github.com/Businge931/sba-api-gateway/internal/api/models"
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
// func (h *OddsHandler) handleOddsRequest(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	serviceCall func(context.Context, interface{}) (interface{}, error),
// 	req interface{},
// ) {
// 	// Read the request body
// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		handleValidationError(w, `{"details": "Failed to read request body"}`)
// 		return
// 	}

// 	// Decode JSON request
// 	if err := json.Unmarshal(body, req); err != nil {
// 		handleValidationError(w, `{"details": "Invalid JSON request"}`)
// 		return
// 	}

// 	// Validate league and date
// 	var league, date string
// 	switch r := req.(type) {
// 	case *models.CreateOddsRequest:
// 		league, date = r.League, r.GameDate
// 	case *models.ReadOddsRequest:
// 		league, date = r.League, r.Date
// 	case *models.UpdateOddsRequest:
// 		league, date = r.League, r.GameDate
// 	case *models.DeleteOddsRequest:
// 		league, date = r.League, r.GameDate
// 	default:
// 		handleValidationError(w, `{"details": "Invalid request type"}`)
// 		return
// 	}

// 	if errMsg, valid := validateLeagueAndDate(league, date); !valid {
// 		handleValidationError(w, errMsg)
// 		return
// 	}

// 	// Convert simple date format to RFC3339 format expected by CRUD service
// 	if parsedDate, err := time.Parse("2006-01-02", date); err == nil {
// 		// Set time to noon UTC to ensure it's a valid future time
// 		rfc3339Date := parsedDate.Format(time.RFC3339)

// 		// Update the request with the formatted date
// 		switch r := req.(type) {
// 		case *models.CreateOddsRequest:
// 			r.GameDate = rfc3339Date
// 			// Convert league name to proper title case for CRUD service
// 			r.League = "English Premier League"
// 		case *models.ReadOddsRequest:
// 			r.Date = rfc3339Date
// 			r.League = "English Premier League"
// 		case *models.UpdateOddsRequest:
// 			r.GameDate = rfc3339Date
// 			r.League = "English Premier League"
// 		case *models.DeleteOddsRequest:
// 			r.GameDate = rfc3339Date
// 			r.League = "English Premier League"
// 		}
// 	}

// 	// Call the service method
// 	res, err := serviceCall(r.Context(), req)
// 	if err != nil {
// 		handleGRPCError(w, err)
// 		return
// 	}

// 	// Write the response
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(res); err != nil {
// 		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
// 		return
// 	}
// }

// CreateOdds handles the creation of odds
func (h *OddsHandler) CreateOdds(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleValidationError(w, "Invalid request body")
		return
	}

	// Validate request
	if _, valid := validateLeagueAndDate(req.League, req.GameDate); !valid {
		handleValidationError(w, "Invalid league or date")
		return
	}

	// Call service
	res, err := h.oddsService.CreateOdds(r.Context(), &req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Odds created successfully",
		"data":    res,
	}); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// ReadOdds handles the reading of odds
func (h *OddsHandler) ReadOdds(w http.ResponseWriter, r *http.Request) {
	var req models.ReadOddsRequest
	if r.Method == http.MethodGet {
		// For GET requests, parse query parameters
		req.League = r.URL.Query().Get("league")
		req.Date = r.URL.Query().Get("date")
	} else {
		// For other methods, parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleValidationError(w, "Invalid request body")
			return
		}
	}

	// Validate request
	if _, valid := validateLeagueAndDate(req.League, req.Date); !valid {
		handleValidationError(w, "Invalid league or date")
		return
	}

	// Call service
	res, err := h.oddsService.ReadOdds(r.Context(), &req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    res,
	}); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// UpdateOdds handles the updating of odds
func (h *OddsHandler) UpdateOdds(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleValidationError(w, "Invalid request body")
		return
	}

	// Validate request
	if _, valid := validateLeagueAndDate(req.League, req.GameDate); !valid {
		handleValidationError(w, "Invalid league or date")
		return
	}

	// Call service
	res, err := h.oddsService.UpdateOdds(r.Context(), &req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Odds updated successfully",
		"data":    res,
	}); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// DeleteOdds handles the deletion of odds
func (h *OddsHandler) DeleteOdds(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteOddsRequest
	if r.Method == http.MethodDelete {
		// For DELETE requests, parse query parameters
		req.League = r.URL.Query().Get("league")
		req.GameDate = r.URL.Query().Get("game_date")
		req.HomeTeam = r.URL.Query().Get("home_team")
		req.AwayTeam = r.URL.Query().Get("away_team")
	} else {
		// For other methods, parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleValidationError(w, "Invalid request body")
			return
		}
	}

	// Validate request
	if _, valid := validateLeagueAndDate(req.League, req.GameDate); !valid {
		handleValidationError(w, "Invalid league or date")
		return
	}

	// Call service
	res, err := h.oddsService.DeleteOdds(r.Context(), &req)
	if err != nil {
		handleGRPCError(w, err)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Odds deleted successfully",
		"data":    res,
	}); err != nil {
		http.Error(w, `{"details": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}
