package utils

import (
	"encoding/json"
	"net/http"

	"github.com/a1mart/kafkaesque/internal/midas/crudmaps/errors"
	"github.com/a1mart/kafkaesque/internal/midas/crudmaps/models"
)

func WriteError(w http.ResponseWriter, code errors.ErrorCode, message string, status int, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := models.ErrorResponse{
		Code:    string(code),
		Message: message,
		Details: details,
	}

	json.NewEncoder(w).Encode(resp)
}
