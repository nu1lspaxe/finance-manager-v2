package utils

import (
	"fmt"
)

const (
	ErrStatusCode = iota + 1
	ErrTokenAlg
	ErrTokenInvalid
)

type BankRecordError struct {
	code    int
	message string
}

func (e *BankRecordError) Error() string {
	return e.message
}

func NewBankRecordError(code int, opts ...string) *BankRecordError {
	message := GetErrorMessage(code, opts...)
	return &BankRecordError{code: code, message: message}
}

func GetErrorMessage(code int, opts ...string) string {
	switch code {
	case ErrStatusCode:
		return fmt.Sprintf("status code error: %v", opts)
	case ErrTokenAlg:
		return fmt.Sprintf("token algorithm error: %v", opts)
	case ErrTokenInvalid:
		return fmt.Sprintf("invalid token: %v", opts)
	default:
		return "unknown error"
	}
}
