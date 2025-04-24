package utils

import (
	"fmt"
)

const (
	ErrUserNotFound = iota + 1
	ErrUserExists
	ErrUserEmailExists
	ErrUserInvalid
	ErrPasswdInvalid
	ErrAccountExists
	ErrStatusCode
	ErrTokenAlg
	ErrTokenInvalid
	ErrUpdateAccountBalance
)

type UserError struct {
	code    int
	message string
}

func (e *UserError) Error() string {
	return e.message
}

func NewUserError(code int, opts ...string) *UserError {
	message := GetErrorMessage(code, opts...)
	return &UserError{code: code, message: message}
}

func GetErrorMessage(code int, opts ...string) string {
	switch code {
	case ErrUserNotFound:
		return fmt.Sprintf("user not found: %v", opts)
	case ErrUserExists:
		return fmt.Sprintf("user already exists: %v", opts)
	case ErrUserEmailExists:
		return fmt.Sprintf("email already exists: %v", opts)
	case ErrUserInvalid:
		return fmt.Sprintf("invalid user: %v", opts)
	case ErrPasswdInvalid:
		return fmt.Sprintf("invalid password: %v", opts)
	case ErrAccountExists:
		return fmt.Sprintf("account already exists: %v", opts)
	case ErrStatusCode:
		return fmt.Sprintf("status code error: %v", opts)
	case ErrTokenAlg:
		return fmt.Sprintf("token algorithm error: %v", opts)
	case ErrTokenInvalid:
		return fmt.Sprintf("invalid token: %v", opts)
	case ErrUpdateAccountBalance:
		return fmt.Sprintf("failed to update account balance: %v", opts)
	default:
		return "unknown error"
	}
}
