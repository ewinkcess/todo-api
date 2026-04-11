package utils

import (
	"net/http"
)

type AppError struct {
	Code    int
	Message string
	Errors  interface{}
}

func (e AppError) Error() string {
	return e.Message
}

// BadRequest error 401
func BadRequest(message string, errors interface{}) AppError {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Errors:  errors,
	}
}

// Unauthorized error 401
func Unauthorized(message string) AppError {
	return AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Errors:  message,
	}
}

// NotFound error 404
func NotFound(message string) AppError {
	return AppError{
		Code:    http.StatusNotFound,
		Message: message,
		Errors:  message,
	}
}

// InternalServerError error 500
func InternalServerError(message string) AppError {
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Errors:  message,
	}
}

// ValidationError error validasi input
func ValidationError(errors interface{}) AppError {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: "Validasi gagal",
		Errors:  errors,
	}
}
