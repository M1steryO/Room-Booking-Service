package helpers

import (
	"github.com/M1steryO/Room-Booking-Service/pkg/logger"
	"net/http"

	"github.com/M1steryO/Room-Booking-Service/internal/delivery/http/models"
	"github.com/M1steryO/Room-Booking-Service/internal/domain"
)

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	appErr := domain.AsAppError(err)
	logger.Error(
		err.Error(),
		"method", r.Method,
		"path", r.URL.Path,
		"error_code", appErr.Code,
		"error_message", appErr.Message,
		"http_status", appErr.HTTPStatus,
	)
	models.WriteError(w, err)
}
