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
)
