package models

type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type VerifyTokenRequest struct {
	Token string `json:"token"`
}
type VerifyTokenResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email"`
}

type RequestPasswordResetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ChangePasswordRequest struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}
