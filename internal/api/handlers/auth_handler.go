package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Businge931/sba-api-gateway/internal/app/domain"
	"github.com/Businge931/sba-api-gateway/internal/app/service"
	"github.com/Businge931/sba-api-gateway/proto"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// handleRequest is a helper function to handle common logic for HTTP requests
func handleRequest[T any, U any](
	w http.ResponseWriter,
	r *http.Request,
	validateFunc func(*T) bool,
	convertFunc func(*T) *U,
	serviceFunc func(context.Context, *U) (*domain.GenericResponse, error),
) {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"details": "Invalid request"}`, http.StatusForbidden)
		return
	}

	// Validate request
	if !validateFunc(&req) {
		http.Error(w, `{"details": "Invalid request data"}`, http.StatusForbidden)
		return
	}

	// Convert proto request to domain request
	domainReq := convertFunc(&req)

	// Call the service method
	res, err := serviceFunc(r.Context(), domainReq)
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

// Login handler
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r,
		ValidateLoginRequest,
		func(req *proto.LoginRequest) *domain.LoginRequest {
			return &domain.LoginRequest{
				Username: req.Username,
				Password: req.Password,
			}
		},
		func(ctx context.Context, req *domain.LoginRequest) (*domain.GenericResponse, error) {
			res, err := h.authService.Login(ctx, req)
			if err != nil {
				// Return nil response with error to trigger proper error handling
				return nil, err
			}
			return &domain.GenericResponse{
				Success: res.Success,
				Message: res.Message,
				Token:   res.Token,
			}, nil
		},
	)
}

// Register handler
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r,
		ValidateRegisterRequest,
		func(req *proto.RegisterRequest) *domain.RegisterRequest {
			return &domain.RegisterRequest{
				Username: req.Username,
				Password: req.Password,
			}
		},
		func(ctx context.Context, req *domain.RegisterRequest) (*domain.GenericResponse, error) {
			res, err := h.authService.Register(ctx, req)
			return &domain.GenericResponse{
				Success: res.Success,
				Message: res.Message,
			}, err
		},
	)
}

// VerifyTokenMiddleware verifies the JWT token before allowing access to protected endpoints
func (h *AuthHandler) VerifyTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from the Authorization header
		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, `{"details": "Missing token"}`, http.StatusUnauthorized)
			return
		}

		// Call AuthService to verify the token
		res, err := h.authService.VerifyToken(r.Context(), token)
		if err != nil {
			handleGRPCError(w, err)
			return
		}

		// Check if the token is valid
		if !res.Success {
			http.Error(w, `{"details": "`+res.Message+`"}`, http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// ValidateLoginRequest validates the LoginRequest fields
func ValidateLoginRequest(req *proto.LoginRequest) bool {
	return req.GetUsername() != "" && req.GetPassword() != ""
}

// ValidateRegisterRequest validates the RegisterRequest fields
func ValidateRegisterRequest(req *proto.RegisterRequest) bool {
	return req.GetUsername() != "" && req.GetPassword() != ""
}

// handleGRPCError maps gRPC errors to HTTP status codes and provides detailed error messages
func handleGRPCError(w http.ResponseWriter, err error) {
	// Set content type for proper JSON response
	w.Header().Set("Content-Type", "application/json")

	st, ok := status.FromError(err)
	if !ok {
		// For non-gRPC errors, use a generic error
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Internal server error",
		})
		return
	}

	// Get the error message from the status
	errorMessage := st.Message()

	// Default HTTP status
	httpStatus := http.StatusInternalServerError

	// Map gRPC status codes to HTTP status codes
	switch st.Code() {
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
	case codes.ResourceExhausted:
		httpStatus = http.StatusTooManyRequests
	case codes.Unavailable:
		httpStatus = http.StatusServiceUnavailable
	}

	// Send error response
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   errorMessage,
	})
}

// writeJSONError writes a JSON formatted error response
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	// Create a properly formatted JSON error response
	response := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}{
		Success: false,
		Message: "",
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}