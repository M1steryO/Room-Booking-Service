package domain

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string { return e.Message }

func NewAppError(code, message string, httpStatus int) error {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func InvalidRequest(message string) error {
	return NewAppError("INVALID_REQUEST", message, http.StatusBadRequest)
}

func Unauthorized(message string) error {
	return NewAppError("UNAUTHORIZED", message, http.StatusUnauthorized)
}

func Forbidden(message string) error { return NewAppError("FORBIDDEN", message, http.StatusForbidden) }

func NotFound(message string) error { return NewAppError("NOT_FOUND", message, http.StatusNotFound) }

func RoomNotFound() error {
	return NewAppError("ROOM_NOT_FOUND", "room not found", http.StatusNotFound)
}

func SlotNotFound() error {
	return NewAppError("SLOT_NOT_FOUND", "slot not found", http.StatusNotFound)
}

func BookingNotFound() error {
	return NewAppError("BOOKING_NOT_FOUND", "booking not found", http.StatusNotFound)
}

func ScheduleExists() error {
	return NewAppError("SCHEDULE_EXISTS", "schedule for this room already exists and cannot be changed", http.StatusConflict)
}

func SlotAlreadyBooked() error {
	return NewAppError("SLOT_ALREADY_BOOKED", "slot is already booked", http.StatusConflict)
}

func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}
}
