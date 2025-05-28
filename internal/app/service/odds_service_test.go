package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOddsClient is a mock implementation of the OddsClient interface
type MockOddsClient struct {
	mock.Mock
}

func (m *MockOddsClient) CreateOdds(ctx context.Context, req *domain.CreateOddsRequest) (*domain.CreateOddsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.CreateOddsResponse), args.Error(1)
}

func (m *MockOddsClient) ReadOdds(ctx context.Context, req *domain.ReadOddsRequest) (*domain.ReadOddsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ReadOddsResponse), args.Error(1)
}

func (m *MockOddsClient) UpdateOdds(ctx context.Context, req *domain.UpdateOddsRequest) (*domain.UpdateOddsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UpdateOddsResponse), args.Error(1)
}

func (m *MockOddsClient) DeleteOdds(ctx context.Context, req *domain.DeleteOddsRequest) (*domain.DeleteOddsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DeleteOddsResponse), args.Error(1)
}

func TestOddsService_CreateOdds(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		mockSetup   func(*MockOddsClient)
		req         *domain.CreateOddsRequest
		want        *domain.CreateOddsResponse
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful odds creation",
			mockSetup: func(moc *MockOddsClient) {
				moc.On("CreateOdds", mock.Anything, &domain.CreateOddsRequest{
					League:          "Premier League",
					HomeTeam:        "Arsenal",
					AwayTeam:        "Chelsea",
					HomeTeamWinOdds: 2.1,
					AwayTeamWinOdds: 3.2,
					DrawOdds:        3.5,
					GameDate:        "2025-05-30",
				}).Return(&domain.CreateOddsResponse{
					Success: true,
					Message: "Odds created successfully",
					Details: "Match between Arsenal and Chelsea on 2025-05-30",
				}, nil)
			},
			req: &domain.CreateOddsRequest{
				League:          "Premier League",
				HomeTeam:        "Arsenal",
				AwayTeam:        "Chelsea",
				HomeTeamWinOdds: 2.1,
				AwayTeamWinOdds: 3.2,
				DrawOdds:        3.5,
				GameDate:        "2025-05-30",
			},
			want: &domain.CreateOddsResponse{
				Success: true,
				Message: "Odds created successfully",
				Details: "Match between Arsenal and Chelsea on 2025-05-30",
			},
			wantErr: false,
		},
		{
			name: "duplicate odds",
			mockSetup: func(moc *MockOddsClient) {
				moc.On("CreateOdds", mock.Anything, mock.Anything).
					Return((*domain.CreateOddsResponse)(nil), errors.New("odds already exists for this match"))
			},
			req: &domain.CreateOddsRequest{
				League:   "Premier League",
				HomeTeam: "Arsenal",
				AwayTeam: "Chelsea",
				GameDate: "2025-05-30",
			},
			want:        nil,
			wantErr:     true,
			expectedErr: errors.New("odds already exists for this match"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockOddsClient)
			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			svc := &oddsService{
				client: mockClient,
			}

			got, err := svc.CreateOdds(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestOddsService_ReadOdds(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		mockSetup   func(*MockOddsClient)
		req         *domain.ReadOddsRequest
		want        *domain.ReadOddsResponse
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful read odds by league and date",
			mockSetup: func(moc *MockOddsClient) {
				moc.On("ReadOdds", mock.Anything, &domain.ReadOddsRequest{
					League: "Premier League",
					Date:   "2025-05-30",
				}).Return(&domain.ReadOddsResponse{
					Odds: []domain.CreateOddsRequest{
						{
							League:          "Premier League",
							HomeTeam:        "Arsenal",
							AwayTeam:        "Chelsea",
							HomeTeamWinOdds: 2.1,
							AwayTeamWinOdds: 3.2,
							DrawOdds:        3.5,
							GameDate:        "2025-05-30",
						},
					},
					Details: "Found 1 matches",
				}, nil)
			},
			req: &domain.ReadOddsRequest{
				League: "Premier League",
				Date:   "2025-05-30",
			},
			want: &domain.ReadOddsResponse{
				Odds: []domain.CreateOddsRequest{
					{
						League:          "Premier League",
						HomeTeam:        "Arsenal",
						AwayTeam:        "Chelsea",
						HomeTeamWinOdds: 2.1,
						AwayTeamWinOdds: 3.2,
						DrawOdds:        3.5,
						GameDate:        "2025-05-30",
					},
				},
				Details: "Found 1 matches",
			},
			wantErr: false,
		},
		{
			name: "no odds found",
			mockSetup: func(moc *MockOddsClient) {
				moc.On("ReadOdds", mock.Anything, mock.Anything).
					Return((*domain.ReadOddsResponse)(nil), errors.New("no odds found for the given criteria"))
			},
			req: &domain.ReadOddsRequest{
				League: "Non-Existent League",
				Date:   "2025-01-01",
			},
			want:        nil,
			wantErr:     true,
			expectedErr: errors.New("no odds found for the given criteria"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockOddsClient)
			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			svc := &oddsService{
				client: mockClient,
			}

			got, err := svc.ReadOdds(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestOddsService_UpdateOdds(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		mockSetup   func(*MockOddsClient)
		req         *domain.UpdateOddsRequest
		want        *domain.UpdateOddsResponse
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful odds update",
			mockSetup: func(moc *MockOddsClient) {
				moc.On("UpdateOdds", mock.Anything, &domain.UpdateOddsRequest{
					League:          "Premier League",
					HomeTeam:        "Arsenal",
					AwayTeam:        "Chelsea",
					HomeTeamWinOdds: 2.2, // Updated odds
					AwayTeamWinOdds: 3.3, // Updated odds
					DrawOdds:        3.6, // Updated odds
					GameDate:        "2025-05-30",
				}).Return(&domain.UpdateOddsResponse{
					Success: true,
					Message: "Odds updated successfully",
				}, nil)
			},
			req: &domain.UpdateOddsRequest{
				League:          "Premier League",
				HomeTeam:        "Arsenal",
				AwayTeam:        "Chelsea",
				HomeTeamWinOdds: 2.2,
				AwayTeamWinOdds: 3.3,
				DrawOdds:        3.6,
				GameDate:        "2025-05-30",
			},
			want: &domain.UpdateOddsResponse{
				Success: true,
				Message: "Odds updated successfully",
			},
			wantErr: false,
		},
		{
			name: "odds not found for update",
			mockSetup: func(moc *MockOddsClient) {
				moc.On("UpdateOdds", mock.Anything, mock.Anything).
					Return((*domain.UpdateOddsResponse)(nil), errors.New("no odds found to update"))
			},
			req: &domain.UpdateOddsRequest{
				League:   "Non-Existent League",
				HomeTeam: "Team A",
				AwayTeam: "Team B",
				GameDate: "2025-01-01",
			},
			want:        nil,
			wantErr:     true,
			expectedErr: errors.New("no odds found to update"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockOddsClient)
			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			svc := &oddsService{
				client: mockClient,
			}

			got, err := svc.UpdateOdds(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestOddsService_DeleteOdds(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		mockSetup   func(*MockOddsClient)
		req         *domain.DeleteOddsRequest
		want        *domain.DeleteOddsResponse
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful odds deletion",
			mockSetup: func(moc *MockOddsClient) {
				moc.On("DeleteOdds", mock.Anything, &domain.DeleteOddsRequest{
					League:   "Premier League",
					HomeTeam: "Arsenal",
					AwayTeam: "Chelsea",
					GameDate: "2025-05-30",
				}).Return(&domain.DeleteOddsResponse{
					Success: true,
					Message: "Odds deleted successfully",
				}, nil)
			},
			req: &domain.DeleteOddsRequest{
				League:   "Premier League",
				HomeTeam: "Arsenal",
				AwayTeam: "Chelsea",
				GameDate: "2025-05-30",
			},
			want: &domain.DeleteOddsResponse{
				Success: true,
				Message: "Odds deleted successfully",
			},
			wantErr: false,
		},
		{
			name: "odds not found for deletion",
			mockSetup: func(moc *MockOddsClient) {
				moc.On("DeleteOdds", mock.Anything, mock.Anything).
					Return((*domain.DeleteOddsResponse)(nil), errors.New("no odds found to delete"))
			},
			req: &domain.DeleteOddsRequest{
				League:   "Non-Existent League",
				HomeTeam: "Team A",
				AwayTeam: "Team B",
				GameDate: "2025-01-01",
			},
			want:        nil,
			wantErr:     true,
			expectedErr: errors.New("no odds found to delete"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockOddsClient)
			if tt.mockSetup != nil {
				tt.mockSetup(mockClient)
			}

			svc := &oddsService{
				client: mockClient,
			}

			got, err := svc.DeleteOdds(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
