package utils

import "fmt"

const (
	ErrTokenAlg = iota + 1
	ErrTokenInvalid
)

type RecordError struct {
	code    int
	message string
}

func NewRecordError(code int, opts ...string) *RecordError {
	message := GetErrorMessage(code, opts...)
	return &RecordError{code: code, message: message}
}

func (e *RecordError) Error() string {
	return e.message
}

func GetErrorMessage(code int, opts ...string) string {
	switch code {
	case ErrTokenAlg:
		return fmt.Sprintf("token algorithm error: %v", opts)
	case ErrTokenInvalid:
		return fmt.Sprintf("invalid token: %v", opts)
	default:
		return "unknown error"
	}
}
