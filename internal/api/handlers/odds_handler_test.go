package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Businge931/sba-api-gateway/internal/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	validToken = "valid-token"
	// The validation expects dates in YYYY-MM-DD format
	// The service expects dates in RFC3339 format
	shortDate   = "2023-01-01"
	rfc3339Date = "2023-01-01T00:00:00Z"
)

// MockOddsService is a mock implementation of OddsHandlerService
type MockOddsService struct {
	mock.Mock
}

func (m *MockOddsService) CreateOdds(ctx context.Context, req *models.CreateOddsRequest) (*models.CreateOddsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.CreateOddsResponse), args.Error(1)
}

func (m *MockOddsService) ReadOdds(ctx context.Context, req *models.ReadOddsRequest) (*models.ReadOddsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.ReadOddsResponse), args.Error(1)
}

func (m *MockOddsService) UpdateOdds(ctx context.Context, req *models.UpdateOddsRequest) (*models.UpdateOddsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.UpdateOddsResponse), args.Error(1)
}

func (m *MockOddsService) DeleteOdds(ctx context.Context, req *models.DeleteOddsRequest) (*models.DeleteOddsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.DeleteOddsResponse), args.Error(1)
}

func createTestRequest(method, url string, body interface{}) *http.Request {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	req.Header.Set("Authorization", validToken)
	return req
}

func TestCreateOdds(t *testing.T) {
	tests := []struct {
		name           string
		setupMock     func(*MockOddsService)
		requestBody   *models.CreateOddsRequest
		wantStatus    int
		wantResponse  string
		validateError func(*testing.T, *httptest.ResponseRecorder, string)
	}{
		{
			name: "successful creation",
			setupMock: func(m *MockOddsService) {
				m.On("CreateOdds", mock.Anything, mock.Anything).Return(&models.CreateOddsResponse{
					Success: true,
					Message: "Odds created successfully",
					Details: "Match odds for Arsenal vs Chelsea created",
				}, nil)
			},
			requestBody: &models.CreateOddsRequest{
				League:          "english premier league",
				GameDate:        shortDate,
				HomeTeam:        "Arsenal",
				AwayTeam:        "Chelsea",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.8,
				DrawOdds:        3.2,
			},
			wantStatus: http.StatusOK,
			validateError: func(t *testing.T, rr *httptest.ResponseRecorder, body string) {
				assert.Contains(t, body, `"success":true`)
				assert.Contains(t, body, `"message":"Odds created successfully"`)
			},
		},
		{
			name: "invalid league",
			setupMock: func(m *MockOddsService) {
				// No expectations as validation fails before service call
			},
			requestBody: &models.CreateOddsRequest{
				League:   "invalid league",
				GameDate: shortDate,
			},
			wantStatus:   http.StatusForbidden,
		},
		{
			name: "invalid date format",
			setupMock: func(m *MockOddsService) {
				// No expectations as validation fails before service call
			},
			requestBody: &models.CreateOddsRequest{
				League:   "english premier league",
				GameDate: "2023/01/01", // Invalid format
			},
			wantStatus:   http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockOddsService{}
			handler := NewOddsHandler(mockSvc)

			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			req := createTestRequest(http.MethodPost, "/odds", tt.requestBody)
			rr := httptest.NewRecorder()

			handler.CreateOdds(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			// Only try to unmarshal as JSON if the content type is application/json
			contentType := rr.Header().Get("Content-Type")
			if strings.Contains(contentType, "application/json") {
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				// Only check success status for successful responses
				if tt.wantStatus == http.StatusOK {
					if success, ok := response["success"].(bool); !ok || !success {
						t.Errorf("handler returned unexpected success status: %v", success)
					}
				}
			}

			if tt.validateError != nil {
				tt.validateError(t, rr, rr.Body.String())
			}
		})
	}
}

func TestReadOdds(t *testing.T) {
	tests := []struct {
		name           string
		setupMock     func(*MockOddsService)
		queryParams   map[string]string
		wantStatus    int
		wantResponse  string
		validateError func(*testing.T, *httptest.ResponseRecorder, string)
	}{
		{
			name: "successful read via GET",
			setupMock: func(m *MockOddsService) {
				m.On("ReadOdds", mock.Anything, mock.MatchedBy(func(req *models.ReadOddsRequest) bool {
					// The service should receive the date in the format it was sent (YYYY-MM-DD)
					return req.League == "english premier league" && req.Date == shortDate
				})).Return(&models.ReadOddsResponse{
					Odds: []models.CreateOddsRequest{
						{
							League:          "english premier league",
							HomeTeam:        "Arsenal",
							AwayTeam:        "Chelsea",
							HomeTeamWinOdds: 2.5,
							AwayTeamWinOdds: 2.8,
							DrawOdds:        3.2,
							GameDate:        rfc3339Date, // The response includes the date in RFC3339 format
						},
					},
					Details: "Successfully retrieved odds",
				}, nil)
			},
			queryParams: map[string]string{
				"league": "english premier league",
				"date":   shortDate, // Validation expects YYYY-MM-DD format
			},
			wantStatus: http.StatusOK,
			validateError: func(t *testing.T, rr *httptest.ResponseRecorder, body string) {
				var response map[string]interface{}
				if err := json.Unmarshal([]byte(body), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				
				// Check the success status
				assert.True(t, response["success"].(bool))
				
				// Check the data field exists and contains the expected odds
				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok, "Expected 'data' field in response")
				
				// Check the odds array
				odds, ok := data["odds"].([]interface{})
				assert.True(t, ok && len(odds) > 0, "Expected 'odds' array in response data")
				
				// Check the first odd in the array
				odd, ok := odds[0].(map[string]interface{})
				assert.True(t, ok, "Expected odd to be an object")
				
				// Check the fields
				assert.Equal(t, "Arsenal", odd["home_team"])
				assert.Equal(t, "Chelsea", odd["away_team"])
				assert.Equal(t, 2.5, odd["home_team_win_odds"])
			},
		},
		{
			name: "missing required parameters",
			setupMock: func(m *MockOddsService) {
				// No expectations as validation should fail before service call
			},
			queryParams: map[string]string{
				// Missing required parameters
			},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockOddsService{}
			handler := NewOddsHandler(mockSvc)

			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Build URL with query parameters
			path := "/odds"
			if len(tt.queryParams) > 0 {
				params := url.Values{}
				for k, v := range tt.queryParams {
					params.Add(k, v)
				}
				path = path + "?" + params.Encode()
			}
			req := createTestRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()

			handler.ReadOdds(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}

			if tt.validateError != nil {
				tt.validateError(t, rr, rr.Body.String())
			}
		})
	}
}

func TestUpdateOdds(t *testing.T) {
	tests := []struct {
		name           string
		setupMock     func(*MockOddsService)
		requestBody   *models.UpdateOddsRequest
		wantStatus    int
		wantResponse  string
		validateError func(*testing.T, *httptest.ResponseRecorder, string)
	}{
		{
			name: "successful update",
			setupMock: func(m *MockOddsService) {
				m.On("UpdateOdds", mock.Anything, mock.Anything).Return(&models.UpdateOddsResponse{
					Success: true,
					Message: "Odds updated successfully",
				}, nil)
			},
			requestBody: &models.UpdateOddsRequest{
				League:          "english premier league",
				GameDate:        shortDate,
				HomeTeam:        "Arsenal",
				AwayTeam:        "Chelsea",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.8,
				DrawOdds:        3.2,
			},
			wantStatus: http.StatusOK,
			validateError: func(t *testing.T, rr *httptest.ResponseRecorder, body string) {
				assert.Contains(t, body, `"success":true`)
				assert.Contains(t, body, `"message":"Odds updated successfully"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockOddsService{}
			handler := NewOddsHandler(mockSvc)

			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/odds", bytes.NewBuffer(body))
			req.Header.Set("Authorization", "Bearer "+validToken)
			w := httptest.NewRecorder()

			handler.UpdateOdds(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantResponse != "" {
				assert.JSONEq(t, tt.wantResponse, w.Body.String())
			}
			if tt.validateError != nil {
				tt.validateError(t, w, w.Body.String())
			}
		})
	}
}

func TestDeleteOdds(t *testing.T) {
	tests := []struct {
		name           string
		setupMock     func(*MockOddsService)
		requestBody   *models.DeleteOddsRequest
		wantStatus    int
		wantResponse  string
		validateError func(*testing.T, *httptest.ResponseRecorder, string)
	}{
		{
			name: "successful deletion",
			setupMock: func(m *MockOddsService) {
				m.On("DeleteOdds", mock.Anything, mock.Anything).Return(&models.DeleteOddsResponse{
					Success: true,
					Message: "Odds deleted successfully",
				}, nil)
			},
			requestBody: &models.DeleteOddsRequest{
				League:   "english premier league",
				GameDate: shortDate,
				HomeTeam: "Arsenal",
				AwayTeam: "Chelsea",
			},
			wantStatus: http.StatusOK,
			validateError: func(t *testing.T, rr *httptest.ResponseRecorder, body string) {
				assert.Contains(t, body, `"success":true`)
				assert.Contains(t, body, `"message":"Odds deleted successfully"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockOddsService{}
			handler := NewOddsHandler(mockSvc)

			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Build URL with query parameters
			params := url.Values{}
			if tt.requestBody != nil {
				if tt.requestBody.League != "" {
					params.Add("league", tt.requestBody.League)
				}
				if tt.requestBody.GameDate != "" {
					params.Add("game_date", tt.requestBody.GameDate)
				}
				if tt.requestBody.HomeTeam != "" {
					params.Add("home_team", tt.requestBody.HomeTeam)
				}
				if tt.requestBody.AwayTeam != "" {
					params.Add("away_team", tt.requestBody.AwayTeam)
				}
			}

			url := "/odds"
			if len(params) > 0 {
				url = url + "?" + params.Encode()
			}

			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Authorization", "Bearer "+validToken)
			w := httptest.NewRecorder()

			handler.DeleteOdds(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantResponse != "" {
				assert.JSONEq(t, tt.wantResponse, w.Body.String())
			}
			if tt.validateError != nil {
				tt.validateError(t, w, w.Body.String())
			}
		})
	}
}
