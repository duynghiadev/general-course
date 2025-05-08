package errors

type ServiceError struct {
	Message string `json:"message"`
}

var (
	ErrEmptyUsername = "Username cannot be empty"
	ErrEmptyPassword = "Password cannot be empty"
	ErrEmptyEmail    = "Email cannot be empty"
	ErrInvalidEmail  = "Invalid email format"
)
