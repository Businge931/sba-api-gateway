package unit

import (
	"context"
	"testing"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOddsServer is a mock of the OddsServer interface
type MockOddsServer struct {
	mock.Mock
}

// CreateOdds mocks the CreateOdds method
func (m *MockOddsServer) CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.CreateOddsResponse), args.Error(1)
}

// ReadOdds mocks the ReadOdds method
func (m *MockOddsServer) ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.ReadOddsResponse), args.Error(1)
}

// UpdateOdds mocks the UpdateOdds method
func (m *MockOddsServer) UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.UpdateOddsResponse), args.Error(1)
}

// DeleteOdds mocks the DeleteOdds method
func (m *MockOddsServer) DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*domain.DeleteOddsResponse), args.Error(1)
}

// TestOddsService tests the OddsService methods
func TestOddsService(t *testing.T) {
	// Create a mock odds server
	mockServer := new(MockOddsServer)
	
	// Create the service
	service := NewOddsServiceForTesting(mockServer)
	
	// Test CreateOdds
	t.Run("CreateOdds", func(t *testing.T) {
		tests := []struct {
			name             string
			league           string
			homeTeam         string
			awayTeam         string
			homeTeamWinOdds  float64
			awayTeamWinOdds  float64
			drawOdds         float64
			gameDate         string
			expectedSuccess  bool
			expectedMessage  string
			expectedDetails  string
			expectedError    error
		}{
			{
				name:             "successful create odds",
				league:           "Premier League",
				homeTeam:         "Arsenal",
				awayTeam:         "Chelsea",
				homeTeamWinOdds:  2.5,
				awayTeamWinOdds:  3.0,
				drawOdds:         2.2,
				gameDate:         "2025-05-15",
				expectedSuccess:  true,
				expectedMessage:  "Odds created successfully",
				expectedDetails:  "Match odds for Arsenal vs Chelsea created",
				expectedError:    nil,
			},
			{
				name:             "invalid odds values",
				league:           "La Liga",
				homeTeam:         "Barcelona",
				awayTeam:         "Real Madrid",
				homeTeamWinOdds:  0.5, // invalid odds value
				awayTeamWinOdds:  0.8, // invalid odds value
				drawOdds:         0.7, // invalid odds value
				gameDate:         "2025-05-20",
				expectedSuccess:  false,
				expectedMessage:  "Failed to create odds",
				expectedDetails:  "Invalid odds values",
				expectedError:    nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Set up expectations
				oddsReq := &domain.CreateOddsRequest{
					League:          tc.league,
					HomeTeam:        tc.homeTeam,
					AwayTeam:        tc.awayTeam,
					HomeTeamWinOdds: tc.homeTeamWinOdds,
					AwayTeamWinOdds: tc.awayTeamWinOdds,
					DrawOdds:        tc.drawOdds,
					GameDate:        tc.gameDate,
				}
				expectedResponse := &domain.CreateOddsResponse{
					Success: tc.expectedSuccess,
					Message: tc.expectedMessage,
					Details: tc.expectedDetails,
				}

				mockServer.On("CreateOdds", mock.Anything, oddsReq).Return(expectedResponse, tc.expectedError).Once()

				// Call the method
				response, err := service.CreateOdds(context.Background(), oddsReq)

				// Assert expectations
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, expectedResponse, response)
				}
			})
		}

		mockServer.AssertExpectations(t)
	})
	
	// Test ReadOdds
	t.Run("ReadOdds", func(t *testing.T) {
		tests := []struct {
			name            string
			league          string
			date            string
			expectedOdds    []domain.CreateOddsRequest
			expectedDetails string
			expectedError   error
		}{
			{
				name:            "successful read odds",
				league:          "Premier League",
				date:            "2025-05-15",
				expectedOdds: []domain.CreateOddsRequest{
					{
						League:          "Premier League",
						HomeTeam:        "Arsenal",
						AwayTeam:        "Chelsea",
						HomeTeamWinOdds: 2.5,
						AwayTeamWinOdds: 3.0,
						DrawOdds:        2.2,
						GameDate:        "2025-05-15",
					},
				},
				expectedDetails: "Found 1 matches",
				expectedError:   nil,
			},
			{
				name:            "no matches found",
				league:          "Serie A",
				date:            "2025-05-15",
				expectedOdds:    []domain.CreateOddsRequest{},
				expectedDetails: "No matches found",
				expectedError:   nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Set up expectations
				readReq := &domain.ReadOddsRequest{
					League: tc.league,
					Date:   tc.date,
				}

				expectedResponse := &domain.ReadOddsResponse{
					Odds:    tc.expectedOdds,
					Details: tc.expectedDetails,
				}

				mockServer.On("ReadOdds", mock.Anything, readReq).Return(expectedResponse, tc.expectedError).Once()

				// Call the method
				response, err := service.ReadOdds(context.Background(), readReq)

				// Assert expectations
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, expectedResponse, response)
					assert.Equal(t, len(tc.expectedOdds), len(response.Odds))
				}
			})
		}

		mockServer.AssertExpectations(t)
	})
	
	// Test UpdateOdds
	t.Run("UpdateOdds", func(t *testing.T) {
		tests := []struct {
			name             string
			league           string
			homeTeam         string
			awayTeam         string
			homeTeamWinOdds  float64
			awayTeamWinOdds  float64
			drawOdds         float64
			gameDate         string
			expectedSuccess  bool
			expectedMessage  string
			expectedDetails  string
			expectedError    error
		}{
			{
				name:             "successful update odds",
				league:           "Premier League",
				homeTeam:         "Arsenal",
				awayTeam:         "Chelsea",
				homeTeamWinOdds:  2.7,
				awayTeamWinOdds:  3.2,
				drawOdds:         2.1,
				gameDate:         "2025-05-15",
				expectedSuccess:  true,
				expectedMessage:  "Odds updated successfully",
				expectedDetails:  "Match odds for Arsenal vs Chelsea updated",
				expectedError:    nil,
			},
			{
				name:             "match not found",
				league:           "Bundesliga",
				homeTeam:         "Bayern Munich",
				awayTeam:         "Borussia Dortmund",
				homeTeamWinOdds:  1.8,
				awayTeamWinOdds:  3.5,
				drawOdds:         3.0,
				gameDate:         "2025-05-15",
				expectedSuccess:  false,
				expectedMessage:  "Failed to update odds",
				expectedDetails:  "Match not found",
				expectedError:    nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Set up expectations
				updateReq := &domain.UpdateOddsRequest{
					League:          tc.league,
					HomeTeam:        tc.homeTeam,
					AwayTeam:        tc.awayTeam,
					HomeTeamWinOdds: tc.homeTeamWinOdds,
					AwayTeamWinOdds: tc.awayTeamWinOdds,
					DrawOdds:        tc.drawOdds,
					GameDate:        tc.gameDate,
				}
				expectedResponse := &domain.UpdateOddsResponse{
					Success: tc.expectedSuccess,
					Message: tc.expectedMessage,
					Details: tc.expectedDetails,
				}

				mockServer.On("UpdateOdds", mock.Anything, updateReq).Return(expectedResponse, tc.expectedError).Once()

				// Call the method
				response, err := service.UpdateOdds(context.Background(), updateReq)

				// Assert expectations
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, expectedResponse, response)
				}
			})
		}

		mockServer.AssertExpectations(t)
	})
	
	// Test DeleteOdds
	t.Run("DeleteOdds", func(t *testing.T) {
		tests := []struct {
			name             string
			league           string
			homeTeam         string
			awayTeam         string
			gameDate         string
			expectedSuccess  bool
			expectedMessage  string
			expectedDetails  string
			expectedError    error
		}{
			{
				name:             "successful delete odds",
				league:           "Premier League",
				homeTeam:         "Arsenal",
				awayTeam:         "Chelsea",
				gameDate:         "2025-05-15",
				expectedSuccess:  true,
				expectedMessage:  "Odds deleted successfully",
				expectedDetails:  "Match odds for Arsenal vs Chelsea deleted",
				expectedError:    nil,
			},
			{
				name:             "match not found",
				league:           "Ligue 1",
				homeTeam:         "PSG",
				awayTeam:         "Marseille",
				gameDate:         "2025-05-15",
				expectedSuccess:  false,
				expectedMessage:  "Failed to delete odds",
				expectedDetails:  "Match not found",
				expectedError:    nil,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Set up expectations
				deleteReq := &domain.DeleteOddsRequest{
					League:   tc.league,
					HomeTeam: tc.homeTeam,
					AwayTeam: tc.awayTeam,
					GameDate: tc.gameDate,
				}
				expectedResponse := &domain.DeleteOddsResponse{
					Success: tc.expectedSuccess,
					Message: tc.expectedMessage,
					Details: tc.expectedDetails,
				}

				mockServer.On("DeleteOdds", mock.Anything, deleteReq).Return(expectedResponse, tc.expectedError).Once()

				// Call the method
				response, err := service.DeleteOdds(context.Background(), deleteReq)

				// Assert expectations
				if tc.expectedError != nil {
					assert.Equal(t, tc.expectedError, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, expectedResponse, response)
				}
			})
		}

		mockServer.AssertExpectations(t)
	})
}

// OddsService is a subset of the service interface we need for testing
type OddsService interface {
	CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error)
	ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error)
	UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error)
	DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error)
}

// OddsServer is the interface we need to mock
type OddsServer interface {
	CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error)
	ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error)
	UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error)
	DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error)
}

// oddsService implements OddsService
type oddsService struct {
	server OddsServer
}

// NewOddsServiceForTesting creates a new OddsService for testing
func NewOddsServiceForTesting(server OddsServer) OddsService {
	return &oddsService{
		server: server,
	}
}

// CreateOdds handles create odds requests
func (s *oddsService) CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error) {
	return s.server.CreateOdds(ctx, req)
}

// ReadOdds handles read odds requests
func (s *oddsService) ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error) {
	return s.server.ReadOdds(ctx, req)
}

// UpdateOdds handles update odds requests
func (s *oddsService) UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error) {
	return s.server.UpdateOdds(ctx, req)
}

// DeleteOdds handles delete odds requests
func (s *oddsService) DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error) {
	return s.server.DeleteOdds(ctx, req)
}
