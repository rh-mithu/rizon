package response

import (
	"encoding/json"
	"errors"
	"github.com/rh-mithu/rizon/backend/internal/apperrors"
	"log"
	"net/http"
)

// JSON sends a JSON response with the given status code and payload
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// Error sends a JSON error response.
// It tries to unwrap the error to find an AppError and use its status code.
func Error(w http.ResponseWriter, err error) {
	log.Println(err)
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		JSON(w, appErr.Code, map[string]string{"error": appErr.Message})
		return
	}

	// Default to 500 Internal Server Error
	JSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
}
