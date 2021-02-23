package errors

import "net/http"

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func InternalServerError(message string)*ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "conflict",
	}
}

func Conflict(message string)*ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusConflict,
		Error:   "conflict",
	}
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

func NotFoundError(message string) *ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}


func NotBadCredentials(message string) *ResponseError {
	return &ResponseError{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "bad_credentials",
	}
}

func StatusDependency(message string) *ResponseError{
	return &ResponseError{
		Message: message,
		Status:  http.StatusFailedDependency,
		Error:   "dependency_failed",
	}
}
