package apperrors

import (
	"fmt"
	"net/http"
)

// AppError is a custom error struct that includes an HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // Internal error, not exposed to user
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap creates a new AppError wrapping an existing error
func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common Errors
var (
	ErrNotFound            = New(http.StatusNotFound, "Resource not found")
	ErrUnauthorized        = New(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = New(http.StatusForbidden, "Forbidden")
	ErrBadRequest          = New(http.StatusBadRequest, "Bad Request")
	ErrInternalServerError = New(http.StatusInternalServerError, "Internal Server Error")
)

// NewNotFound creates a 404 error with a specific message
func NewNotFound(item string) *AppError {
	return New(http.StatusNotFound, fmt.Sprintf("%s not found", item))
}

// NewBadRequest creates a 400 error with a specific message
func NewBadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

// NewInternalError wraps an internal error with a 500 status
func NewInternalError(err error) *AppError {
	return Wrap(err, http.StatusInternalServerError, "Internal Server Error")
}
