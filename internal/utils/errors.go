package utils

import (
	"fmt"
	"net/http"
)

// AppError represents an application-specific error
type AppError struct {
	Message    string
	StatusCode int
	Cause      error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(message string, statusCode int, cause error) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
		Cause:      cause,
	}
}

// Common error constructors
func NewBadRequestError(message string) *AppError {
	return NewAppError(message, http.StatusBadRequest, nil)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(message, http.StatusNotFound, nil)
}

func NewInternalServerError(message string, cause error) *AppError {
	return NewAppError(message, http.StatusInternalServerError, cause)
}

// ErrorResponse represents the structure of error responses sent to clients
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteErrorResponse writes an error response to the http.ResponseWriter
func WriteErrorResponse(w http.ResponseWriter, err error) {
	appErr, ok := err.(*AppError)
	if !ok {
		appErr = NewInternalServerError("An unexpected error occurred", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: appErr.Message})
}