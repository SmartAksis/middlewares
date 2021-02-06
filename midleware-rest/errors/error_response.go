package errors

import "net/http"

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func MethodForbidden(message string)*ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusForbidden,
		Error:   "forbidden",
	}
}

func NotAuthorized(message string) *ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

func NewBadRequestError(message string) *ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusBadGateway,
		Error:   "bad_request",
	}
}

func NewNotFoundError(message string) *ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}
