package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Status  int
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func New(code string, status int, message string) *AppError {
	return &AppError{Code: code, Status: status, Message: message}
}

func Wrap(code string, status int, message string, cause error) *AppError {
	return &AppError{Code: code, Status: status, Message: message, Cause: cause}
}
