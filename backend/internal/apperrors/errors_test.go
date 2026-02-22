package apperrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New(http.StatusBadRequest, "Bad Request")

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "Bad Request", err.Message)
	assert.Nil(t, err.Err)
}

func TestAppError_Error(t *testing.T) {
	err := New(http.StatusNotFound, "Not Found")
	assert.Equal(t, "Not Found", err.Error())
}

func TestAppError_Unwrap(t *testing.T) {
	innerErr := errors.New("database error")
	appErr := &AppError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Error",
		Err:     innerErr,
	}

	unwrapped := appErr.Unwrap()
	assert.Equal(t, innerErr, unwrapped)
}

func TestWrap(t *testing.T) {
	innerErr := errors.New("database error")
	appErr := Wrap(innerErr, http.StatusInternalServerError, "Something went wrong")

	assert.NotNil(t, appErr)
	assert.Equal(t, http.StatusInternalServerError, appErr.Code)
	assert.Equal(t, "Something went wrong", appErr.Message)
	assert.Equal(t, innerErr, appErr.Err)
}

func TestNewNotFound(t *testing.T) {
	err := NewNotFound("User")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Code)
	assert.Equal(t, "User not found", err.Message)
}

func TestNewBadRequest(t *testing.T) {
	err := NewBadRequest("Invalid input")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, "Invalid input", err.Message)
}

func TestNewInternalError(t *testing.T) {
	innerErr := errors.New("db down")
	err := NewInternalError(innerErr)

	assert.NotNil(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Equal(t, "Internal Server Error", err.Message)
	assert.Equal(t, innerErr, err.Err)
}
