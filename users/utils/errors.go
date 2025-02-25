package utils

const (
	USER_NOT_FOUND    = 1
	USER_EXISTS       = 2
	USER_EMAIL_EXISTS = 3
	USER_INVALID      = 4
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
