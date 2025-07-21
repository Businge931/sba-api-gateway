package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Businge931/sba-api-gateway/internal/api"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
	serverGrpc "github.com/Businge931/sba-api-gateway/internal/server/grpc"
	"github.com/Businge931/sba-api-gateway/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// bufDialer implements a virtual gRPC connection for testing
func bufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}
}

// MockAuthServer implements the proto.AuthServiceServer interface
type MockAuthServer struct {
	proto.UnimplementedAuthServiceServer
}

// Login mocks the login method
func (s *MockAuthServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if req.Email == "testuser" && req.Password == "testpass" {
		return &proto.LoginResponse{
			Success: true,
			Token:   "test-jwt-token",
			Message: "Login successful",
		}, nil
	}
	return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
}

// Register mocks the register method
func (s *MockAuthServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if req.Email == "newuser" {
		return &proto.RegisterResponse{
			Success: true,
			Message: "User registered successfully",
		}, nil
	}
	return nil, status.Error(codes.AlreadyExists, "User already exists")
}

// VerifyToken mocks the verify token method
func (s *MockAuthServer) VerifyToken(ctx context.Context, req *proto.VerifyTokenRequest) (*proto.VerifyTokenResponse, error) {
	if req.Token == "test-jwt-token" {
		return &proto.VerifyTokenResponse{
			Success: true,
			Message: "Token is valid",
		}, nil
	}
	return nil, status.Error(codes.Unauthenticated, "Invalid token")
}

// MockOddsServer implements the proto.OddsServiceServer interface
type MockOddsServer struct {
	proto.UnimplementedOddsServiceServer
}

// CreateOdds mocks the create odds method
func (s *MockOddsServer) CreateOdds(ctx context.Context, req *proto.CreateOddsRequest) (*proto.CreateOddsResponse, error) {
	return &proto.CreateOddsResponse{
		Success: true,
		Message: "Odds created successfully",
		Details: fmt.Sprintf("Match odds for %s vs %s created", req.HomeTeam, req.AwayTeam),
	}, nil
}

// ReadOdds mocks the read odds method
func (s *MockOddsServer) ReadOdds(ctx context.Context, req *proto.ReadOddsRequest) (*proto.ReadOddsResponse, error) {
	odds := []*proto.CreateOddsRequest{
		{
			League:          req.League,
			HomeTeam:        "Arsenal",
			AwayTeam:        "Chelsea",
			HomeTeamWinOdds: 2.5,
			AwayTeamWinOdds: 3.0,
			DrawOdds:        2.2,
			GameDate:        req.Date,
		},
	}

	return &proto.ReadOddsResponse{
		Odds:    odds,
		Details: "Found 1 matches",
	}, nil
}

// UpdateOdds mocks the update odds method
func (s *MockOddsServer) UpdateOdds(ctx context.Context, req *proto.UpdateOddsRequest) (*proto.UpdateOddsResponse, error) {
	return &proto.UpdateOddsResponse{
		Success: true,
		Message: "Odds updated successfully",
		Details: fmt.Sprintf("Match odds for %s vs %s updated", req.HomeTeam, req.AwayTeam),
	}, nil
}

// DeleteOdds mocks the delete odds method
func (s *MockOddsServer) DeleteOdds(ctx context.Context, req *proto.DeleteOddsRequest) (*proto.DeleteOddsResponse, error) {
	return &proto.DeleteOddsResponse{
		Success: true,
		Message: "Odds deleted successfully",
		Details: fmt.Sprintf("Match odds for %s vs %s deleted", req.HomeTeam, req.AwayTeam),
	}, nil
}

// This function sets up a mock gRPC server for Auth service
func setupMockAuthServer(t *testing.T) (*bufconn.Listener, *MockAuthServer) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	mockServer := &MockAuthServer{}
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, mockServer)

	go func() {
		if err := server.Serve(listener); err != nil {
			require.NoError(t, err, "Failed to serve mock auth server")
		}
	}()

	return listener, mockServer
}

// This function sets up a mock gRPC server for Odds service
func setupMockOddsServer(t *testing.T) (*bufconn.Listener, *MockOddsServer) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	mockServer := &MockOddsServer{}
	server := grpc.NewServer()
	proto.RegisterOddsServiceServer(server, mockServer)

	go func() {
		if err := server.Serve(listener); err != nil {
			require.NoError(t, err, "Failed to serve mock odds server")
		}
	}()

	return listener, mockServer
}

// TestAPIGatewayBufconn tests the API Gateway integration with in-memory mocked services
func TestAPIGatewayBufconn(t *testing.T) {
	// Set up mock grpc servers
	authListener, _ := setupMockAuthServer(t)
	oddsListener, _ := setupMockOddsServer(t)

	// Create auth client connection
	authCtx := context.Background()
	authConn, err := grpc.DialContext(
		authCtx,
		"bufnet",
		grpc.WithContextDialer(bufDialer(authListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer authConn.Close()

	// Create odds client connection
	oddsCtx := context.Background()
	oddsConn, err := grpc.DialContext(
		oddsCtx,
		"bufnet",
		grpc.WithContextDialer(bufDialer(oddsListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	defer oddsConn.Close()

	// Initialize gRPC clients
	authGrpcServer := serverGrpc.NewAuthServer(authConn)
	oddsGrpcServer := serverGrpc.NewOddsServer(oddsConn)

	// Initialize services
	authService := service.NewAuthService(authGrpcServer)
	oddsService := service.NewOddsService(oddsGrpcServer)

	// Setup routes
	router := api.SetupRoutes(authService, oddsService)

	t.Run("Login", func(t *testing.T) {
		tests := []struct {
			name            string
			username        string
			password        string
			expectedStatus  int
			expectedSuccess bool
			expectedToken   string
			expectedMessage string
		}{
			{
				name:            "successful login",
				username:        "testuser",
				password:        "testpass",
				expectedStatus:  http.StatusOK,
				expectedSuccess: true,
				expectedToken:   "test-jwt-token",
				expectedMessage: "Login successful",
			},
			{
				name:            "invalid credentials",
				username:        "wronguser",
				password:        "wrongpass",
				expectedStatus:  http.StatusUnauthorized,
				expectedSuccess: false,
				expectedToken:   "",
				expectedMessage: "",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Create login request
				loginReq := map[string]string{
					"username": tc.username,
					"password": tc.password,
				}
				payload, err := json.Marshal(loginReq)
				require.NoError(t, err)

				// Create request
				req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				// Execute request
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				// Check status code
				assert.Equal(t, tc.expectedStatus, rr.Code)

				// For successful responses, validate the response content
				if tc.expectedStatus == http.StatusOK {
					var response map[string]any
					err = json.Unmarshal(rr.Body.Bytes(), &response)
					require.NoError(t, err)

					// Validate response
					assert.Equal(t, tc.expectedSuccess, response["success"])
					assert.Equal(t, tc.expectedToken, response["token"])
				}
			})
		}
	})

	t.Run("CreateOdds", func(t *testing.T) {
		tests := []struct {
			name            string
			league          string
			homeTeam        string
			awayTeam        string
			homeTeamWinOdds float64
			awayTeamWinOdds float64
			drawOdds        float64
			gameDate        string
			token           string
			expectedStatus  int
			expectedSuccess bool
			expectedMessage string
		}{
			{
				name:            "successful create odds",
				league:          "english premier league",
				homeTeam:        "Arsenal",
				awayTeam:        "Chelsea",
				homeTeamWinOdds: 2.5,
				awayTeamWinOdds: 3.0,
				drawOdds:        2.2,
				gameDate:        "2025-05-15",
				token:           "test-jwt-token",
				expectedStatus:  http.StatusOK,
				expectedSuccess: true,
				expectedMessage: "Odds created successfully",
			},
			{
				name:            "unauthorized request",
				league:          "english premier league",
				homeTeam:        "Arsenal",
				awayTeam:        "Chelsea",
				homeTeamWinOdds: 2.5,
				awayTeamWinOdds: 3.0,
				drawOdds:        2.2,
				gameDate:        "2025-05-15",
				token:           "invalid-token",
				expectedStatus:  http.StatusUnauthorized,
				expectedSuccess: false,
				expectedMessage: "",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Create odds request
				createOddsReq := map[string]interface{}{
					"league":             tc.league,
					"home_team":          tc.homeTeam,
					"away_team":          tc.awayTeam,
					"home_team_win_odds": tc.homeTeamWinOdds,
					"away_team_win_odds": tc.awayTeamWinOdds,
					"draw_odds":          tc.drawOdds,
					"game_date":          tc.gameDate,
				}
				payload, err := json.Marshal(createOddsReq)
				require.NoError(t, err)

				// Create request
				req, err := http.NewRequest("POST", "/api/odds/create", bytes.NewBuffer(payload))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", tc.token)

				// Execute request
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				// Check status code
				assert.Equal(t, tc.expectedStatus, rr.Code)

				// For successful responses, validate the response content
				if tc.expectedStatus == http.StatusOK {
					var response map[string]interface{}
					err = json.Unmarshal(rr.Body.Bytes(), &response)
					require.NoError(t, err)

					// Validate response
					assert.Equal(t, tc.expectedSuccess, response["success"])
					assert.Equal(t, tc.expectedMessage, response["message"])
				}
			})
		}
	})

	t.Run("ReadOdds", func(t *testing.T) {
		tests := []struct {
			name           string
			league         string
			date           string
			token          string
			expectedStatus int
		}{
			{
				name:           "successful read odds",
				league:         "english premier league",
				date:           "2025-05-15",
				token:          "test-jwt-token",
				expectedStatus: http.StatusOK,
			},
			{
				name:           "unauthorized request",
				league:         "english premier league",
				date:           "2025-05-15",
				token:          "invalid-token",
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// Create request URL with query parameters
				url := fmt.Sprintf("/api/odds/read?league=%s&date=%s", tc.league, tc.date)

				// Create request
				req, err := http.NewRequest("GET", url, nil)
				require.NoError(t, err)
				req.Header.Set("Authorization", tc.token)

				// Execute request
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				// Check status code
				assert.Equal(t, tc.expectedStatus, rr.Code)

				// For successful responses, validate that a response was received
				if tc.expectedStatus == http.StatusOK {
					assert.NotEmpty(t, rr.Body.String())
				}
			})
		}
	})
}
