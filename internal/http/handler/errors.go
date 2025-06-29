package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// AppError represents application-specific errors with types
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return "unknown error"
}

// ErrorType defines different categories of errors
type ErrorType int

const (
	ErrorTypeValidation ErrorType = iota
	ErrorTypeBadRequest
	ErrorTypeNotFound
	ErrorTypeConflict
	ErrorTypeInternal
)

// NewValidationError creates a validation error
func NewValidationError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Cause:   cause,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Cause:   cause,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Cause:   cause,
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Cause:   cause,
	}
}

// NewInternalError creates an internal server error
func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Cause:   cause,
	}
}

// HandleError processes errors and returns appropriate HTTP responses
func HandleError(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	// Check if it's our custom AppError
	if appErr, ok := err.(*AppError); ok {
		return handleAppError(c, appErr)
	}

	// Handle common error patterns from dependencies
	return handleGenericError(c, err)
}

func handleAppError(c echo.Context, appErr *AppError) error {
	log.Error().Err(appErr).Msgf("handleAppError: processing error type %d with message: %s", appErr.Type, appErr.Message)

	var statusCode int
	var errorCode string

	switch appErr.Type {
	case ErrorTypeValidation:
		statusCode = http.StatusBadRequest
		errorCode = "VALIDATION_ERROR"
	case ErrorTypeNotFound:
		statusCode = http.StatusNotFound
		errorCode = "NOT_FOUND"
	case ErrorTypeConflict:
		statusCode = http.StatusConflict
		errorCode = "CONFLICT"
	case ErrorTypeBadRequest:
		statusCode = http.StatusBadRequest
		errorCode = "BAD_REQUEST"
	default:
		statusCode = http.StatusInternalServerError
		errorCode = "INTERNAL_ERROR"
	}

	response := ErrorResponse{
		Error:   appErr.Error(),
		Code:    errorCode,
		Message: appErr.Message,
	}

	log.Info().Msgf("handleAppError: returning status %d with response: %+v", statusCode, response)
	return c.JSON(statusCode, response)
}

func handleGenericError(c echo.Context, err error) error {
	errMsg := err.Error()
	log.Error().Err(err).Msgf("handleGenericError: processing error message: %s", errMsg)

	// Check for common database errors
	if strings.Contains(errMsg, "duplicate key") || strings.Contains(errMsg, "unique constraint") {
		response := ErrorResponse{
			Error:   "Resource already exists",
			Code:    "CONFLICT",
			Message: "A resource with this information already exists",
		}
		log.Error().Err(err).Msg("handleGenericError: detected duplicate key error, returning conflict")
		return c.JSON(http.StatusConflict, response)
	}

	if strings.Contains(errMsg, "no rows in result set") || strings.Contains(errMsg, "not found") {
		response := ErrorResponse{
			Error:   "Resource not found",
			Code:    "NOT_FOUND",
			Message: "The requested resource was not found",
		}
		log.Error().Err(err).Msg("handleGenericError: detected not found error")
		return c.JSON(http.StatusNotFound, response)
	}

	// Check for UUID parsing errors
	if strings.Contains(errMsg, "invalid UUID") || strings.Contains(errMsg, "uuid: incorrect UUID length") {
		response := ErrorResponse{
			Error:   "Invalid UUID format",
			Code:    "BAD_REQUEST",
			Message: "The provided ID is not a valid UUID",
		}
		log.Error().Err(err).Msg("handleGenericError: detected UUID parsing error")
		return c.JSON(http.StatusBadRequest, response)
	}

	// Check for JSON binding errors
	if strings.Contains(errMsg, "json") || strings.Contains(errMsg, "character") || strings.Contains(errMsg, "syntax") {
		response := ErrorResponse{
			Error:   "Invalid JSON format",
			Code:    "BAD_REQUEST",
			Message: "The request body contains invalid JSON",
		}
		log.Error().Err(err).Msg("handleGenericError: detected JSON binding error")
		return c.JSON(http.StatusBadRequest, response)
	}

	// Default to internal server error
	response := ErrorResponse{
		Error:   "Internal server error",
		Code:    "INTERNAL_ERROR",
		Message: "An unexpected error occurred",
	}
	log.Error().Err(err).Msg("handleGenericError: defaulting to internal server error")
	return c.JSON(http.StatusInternalServerError, response)
}
