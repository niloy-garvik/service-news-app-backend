package utils

import (
	"fmt"
	"net/http"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	StatusCode string      `json:"statusCode"`
	ErrorCode  string      `json:"errorCode"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	StatusCode string      `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// SendErrorResponse sends a JSON formatted error response
func SendErrorResponse(w http.ResponseWriter, statusCode int, errorCode, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := ErrorResponse{
		StatusCode: fmt.Sprintf("%d", statusCode),
		ErrorCode:  errorCode,
		Message:    message,
		Details:    details,
	}

	apiResponse := ConvertToJson(errorResponse)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(apiResponse))
}

// SendSuccessResponse sends a JSON formatted success response
func SendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	successResponse := SuccessResponse{
		StatusCode: fmt.Sprintf("%d", statusCode),
		Message:    message,
		Data:       data,
	}
	apiResponse := ConvertToJson(successResponse)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(apiResponse))
}
