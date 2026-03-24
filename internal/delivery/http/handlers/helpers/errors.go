package helpers

import (
	"github.com/avito-internships/test-backend-1-M1steryO/pkg/logger"
	"net/http"

	"github.com/avito-internships/test-backend-1-M1steryO/internal/delivery/http/models"
	"github.com/avito-internships/test-backend-1-M1steryO/internal/domain"
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
