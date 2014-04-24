package exceptions

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

// login error
type PageNotFoundError struct {
	Message string
	Reason  string
}

func (e *PageNotFoundError) Error() string { return e.Message }

func NewPageNotFoundError(message string) *LoginError { return &LoginError{Message: message} }
