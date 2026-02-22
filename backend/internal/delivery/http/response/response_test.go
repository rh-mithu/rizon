package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"finance-tracker-backend/internal/apperrors"
	"finance-tracker-backend/internal/delivery/http/response" // Assuming package is `response`

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"foo": "bar"}
	response.JSON(w, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var body map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, "bar", body["foo"])
}

func TestJSON_NoPayload(t *testing.T) {
	w := httptest.NewRecorder()
	response.JSON(w, http.StatusAccepted, nil)

	assert.Equal(t, http.StatusAccepted, w.Code)
	// Body should be empty (or newline from encoder) vs completely empty depending on implementation
	// The current implementation:
	// if payload != nil { json.NewEncoder(w).Encode(payload) }
	// So body should be empty.
	assert.Empty(t, w.Body.Bytes())
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("some system error")
	response.Error(w, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestError_AppError(t *testing.T) {
	w := httptest.NewRecorder()
	appErr := apperrors.New(http.StatusBadRequest, "bad request")
	response.Error(w, appErr)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	assert.Equal(t, "bad request", body["error"])
}
