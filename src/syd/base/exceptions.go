package base

// ___________________________________________________________________________

// permisstion denied error
type AccessDeniedError struct {
	Message string
	Reason  string
}

func (e *AccessDeniedError) Error() string { return e.Message }

// ___________________________________________________________________________

// login error
type LoginError struct {
	Message string
	Reason  string
}

func (e *LoginError) Error() string { return e.Message }

func NewLoginError(message string) *LoginError { return &LoginError{Message: message} }

// ___________________________________________________________________________
