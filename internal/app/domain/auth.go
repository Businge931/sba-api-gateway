package domain

type (
	GenericResponse struct {
		Success bool
		Message string
		Token   string
	}

	LoginRequest struct {
		Email    string
		Password string
	}

	LoginResponse struct {
		Success bool
		Token   string
		Message string
	}

	RegisterRequest struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
	}

	RegisterResponse struct {
		Success bool
		Message string
	}

	VerifyTokenRequest struct {
		Token string
	}
	VerifyTokenResponse struct {
		Success bool
		Message string
	}

	RequestPasswordResetRequest struct {
		Email string
	}

	RequestPasswordResetResponse struct {
		Success bool
		Message string
	}

	ChangePasswordRequest struct {
		UserID      string
		OldPassword string
		NewPassword string
	}

	ChangePasswordResponse struct {
		Success bool
		Message string
	}

	ResetPasswordRequest struct {
		Token       string
		NewPassword string
	}

	ResetPasswordResponse struct {
		Success bool
		Message string
		UserID  string
	}
)
