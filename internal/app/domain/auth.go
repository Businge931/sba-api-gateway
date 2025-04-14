package domain

type GenericResponse struct {
	Success bool
	Message string
	Token   string
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Success bool
	Token   string
	Message string
}

type RegisterRequest struct {
	Username string
	Password string
}

type RegisterResponse struct {
	Success bool
	Message string
}

type VerifyTokenRequest struct {
	Token string
}
type VerifyTokenResponse struct {
	Success bool
	Message string
}
