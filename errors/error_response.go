package errors

import "net/http"

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewBadRequestError(message string) *ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusBadGateway,
		Error:   "bad_request",
	}
}
