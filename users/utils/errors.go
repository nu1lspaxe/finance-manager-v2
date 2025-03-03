package utils

const (
	ErrUserNotFound = iota + 1
	ErrUserExists
	ErrUserEmailExists
	ErrUserInvalid
)

type UserError struct {
	code    int
	message string
}

func (e *UserError) Error() string {
	return e.message
}

func NewUserError(code int, message string) *UserError {
	return &UserError{code: code, message: message}
}
