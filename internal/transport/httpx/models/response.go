package models

import (
	"encoding/json"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
	"net/http"
)

type ErrorEnvelope struct {
	Error ErrorDTO `json:"error"`
}

type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, err error) {
	appErr := domain.AsAppError(err)
	WriteJSON(w, appErr.HTTPStatus, ErrorEnvelope{
		Error: ErrorDTO{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}
