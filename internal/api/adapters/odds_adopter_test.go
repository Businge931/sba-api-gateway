package adapters

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Businge931/sba-api-gateway/internal/api/models"
	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOddsService is a mock implementation of the OddsService interface
type MockOddsService struct {
	mock.Mock
}

func (m *MockOddsService) CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CreateOddsResponse), args.Error(1)
}

func (m *MockOddsService) ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ReadOddsResponse), args.Error(1)
}

func (m *MockOddsService) UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UpdateOddsResponse), args.Error(1)
}

func (m *MockOddsService) DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DeleteOddsResponse), args.Error(1)
}

func TestOddsServiceAdapter_CreateOdds(t *testing.T) {
	// Test data
	testTime := time.Date(2023, 5, 28, 19, 0, 0, 0, time.UTC).Format(time.RFC3339)

	// Define test cases
	tests := []struct {
		name           string
		setupMock     func(*MockOddsService)
		req           *models.CreateOddsRequest
		expected      *models.CreateOddsResponse
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful creation",
			setupMock: func(m *MockOddsService) {
				m.On("CreateOdds", mock.Anything, &domain.CreateOddsRequest{
					League:          "EPL",
					GameDate:        testTime,
					HomeTeam:        "Arsenal",
					AwayTeam:        "Chelsea",
					HomeTeamWinOdds: 2.1,
					AwayTeamWinOdds: 3.2,
					DrawOdds:        3.5,
				}).Return(&domain.CreateOddsResponse{
					Success: true,
					Message: "Odds created successfully",
				}, nil)
			},
			req: &models.CreateOddsRequest{
				League:          "EPL",
				GameDate:        testTime,
				HomeTeam:        "Arsenal",
				AwayTeam:        "Chelsea",
				HomeTeamWinOdds: 2.1,
				AwayTeamWinOdds: 3.2,
				DrawOdds:        3.5,
			},
			expected: &models.CreateOddsResponse{
				Success: true,
				Message: "Odds created successfully",
			},
			expectedError: false,
		},
		{
			name: "service error",
			setupMock: func(m *MockOddsService) {
				m.On("CreateOdds", mock.Anything, mock.Anything).Return(
					(*domain.CreateOddsResponse)(nil),
					errors.New("database error"),
				)
			},
			req: &models.CreateOddsRequest{
				League:   "EPL",
				GameDate: testTime,
				HomeTeam: "Arsenal",
				AwayTeam: "Chelsea",
			},
			expected:      nil,
			expectedError: true,
			expectedErrMsg: "database error",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockOddsService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Create adapter with mock service
			adapter := NewOddsServiceAdapter(mockSvc)
			// Call the method under test
			result, err := adapter.CreateOdds(context.Background(), tt.req)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestOddsServiceAdapter_ReadOdds(t *testing.T) {
	// Test data
	testDate := time.Date(2023, 5, 28, 0, 0, 0, 0, time.UTC)
	testDateStr := testDate.Format("2006-01-02")
	gameTime := testDate.Add(2 * time.Hour).Format(time.RFC3339)

	// Define test cases
	tests := []struct {
		name           string
		setupMock     func(*MockOddsService)
		req           *models.ReadOddsRequest
		expected      *models.ReadOddsResponse
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful read",
			setupMock: func(m *MockOddsService) {
				m.On("ReadOdds", mock.Anything, &domain.ReadOddsRequest{
					League: "EPL",
					Date:   testDateStr,
				}).Return(&domain.ReadOddsResponse{
					Odds: []domain.CreateOddsRequest{
						{
							League:          "EPL",
							GameDate:        gameTime,
							HomeTeam:        "Arsenal",
							AwayTeam:        "Chelsea",
							HomeTeamWinOdds: 2.1,
							AwayTeamWinOdds: 3.2,
							DrawOdds:        3.5,
						},
					},
					Details: "2 matches found",
				}, nil)
			},
			req: &models.ReadOddsRequest{
				League: "EPL",
				Date:   testDateStr,
			},
			expected: &models.ReadOddsResponse{
				Odds: []models.CreateOddsRequest{
					{
						League:          "EPL",
						GameDate:        gameTime,
						HomeTeam:        "Arsenal",
						AwayTeam:        "Chelsea",
						HomeTeamWinOdds: 2.1,
						AwayTeamWinOdds: 3.2,
						DrawOdds:        3.5,
					},
				},
				Details: "2 matches found",
			},
			expectedError: false,
		},
		{
			name: "not found",
			setupMock: func(m *MockOddsService) {
				m.On("ReadOdds", mock.Anything, &domain.ReadOddsRequest{
					League: "LaLiga",
					Date:   testDateStr,
				}).Return((*domain.ReadOddsResponse)(nil), errors.New("no odds found"))
			},
			req: &models.ReadOddsRequest{
				League: "LaLiga",
				Date:   testDateStr,
			},
			expected:      nil,
			expectedError: true,
			expectedErrMsg: "no odds found",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockOddsService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Create adapter with mock service
			adapter := NewOddsServiceAdapter(mockSvc)
			// Call the method under test
			result, err := adapter.ReadOdds(context.Background(), tt.req)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestOddsServiceAdapter_UpdateOdds(t *testing.T) {
	// Test data
	testTime := time.Date(2023, 5, 28, 19, 0, 0, 0, time.UTC).Format(time.RFC3339)

	// Define test cases
	tests := []struct {
		name           string
		setupMock     func(*MockOddsService)
		req           *models.UpdateOddsRequest
		expected      *models.UpdateOddsResponse
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful update",
			setupMock: func(m *MockOddsService) {
				m.On("UpdateOdds", mock.Anything, &domain.UpdateOddsRequest{
					League:          "EPL",
					GameDate:        testTime,
					HomeTeam:        "Arsenal",
					AwayTeam:        "Chelsea",
					HomeTeamWinOdds: 2.2,
					AwayTeamWinOdds: 3.3,
					DrawOdds:        3.6,
				}).Return(&domain.UpdateOddsResponse{
					Success: true,
					Message: "Odds updated successfully",
				}, nil)
			},
			req: &models.UpdateOddsRequest{
				League:          "EPL",
				GameDate:        testTime,
				HomeTeam:        "Arsenal",
				AwayTeam:        "Chelsea",
				HomeTeamWinOdds: 2.2,
				AwayTeamWinOdds: 3.3,
				DrawOdds:        3.6,
			},
			expected: &models.UpdateOddsResponse{
				Success: true,
				Message: "Odds updated successfully",
			},
			expectedError: false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockOddsService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Create adapter with mock service
			adapter := NewOddsServiceAdapter(mockSvc)


			// Call the method under test
			result, err := adapter.UpdateOdds(context.Background(), tt.req)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestOddsServiceAdapter_DeleteOdds(t *testing.T) {
	// Test data
	testTime := time.Date(2023, 5, 28, 19, 0, 0, 0, time.UTC).Format(time.RFC3339)

	// Define test cases
	tests := []struct {
		name           string
		setupMock     func(*MockOddsService)
		req           *models.DeleteOddsRequest
		expected      *models.DeleteOddsResponse
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful deletion",
			setupMock: func(m *MockOddsService) {
				m.On("DeleteOdds", mock.Anything, &domain.DeleteOddsRequest{
					League:   "EPL",
					GameDate: testTime,
				}).Return(&domain.DeleteOddsResponse{
					Success: true,
					Message: "Odds deleted successfully",
				}, nil)
			},
			req: &models.DeleteOddsRequest{
				League:   "EPL",
				GameDate: testTime,
			},
			expected: &models.DeleteOddsResponse{
				Success: true,
				Message: "Odds deleted successfully",
			},
			expectedError: false,
		},
		{
			name: "not found",
			setupMock: func(m *MockOddsService) {
				m.On("DeleteOdds", mock.Anything, &domain.DeleteOddsRequest{
					League:   "LaLiga",
					GameDate: testTime,
				}).Return((*domain.DeleteOddsResponse)(nil), errors.New("odds not found"))
			},
			req: &models.DeleteOddsRequest{
				League:   "LaLiga",
				GameDate: testTime,
			},
			expected:      nil,
			expectedError: true,
			expectedErrMsg: "odds not found",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockSvc := &MockOddsService{}
			if tt.setupMock != nil {
				tt.setupMock(mockSvc)
			}

			// Create adapter with mock service
			adapter := NewOddsServiceAdapter(mockSvc)
			// Call the method under test
			result, err := adapter.DeleteOdds(context.Background(), tt.req)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Verify mock expectations
			mockSvc.AssertExpectations(t)
		})
	}
}
