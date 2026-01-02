package utils

import "net/http"

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

func NewBadRequest(msg string) *APIError {
	return &APIError{Status: http.StatusBadRequest, Message: msg}
}

func NewConflict(msg string) *APIError {
	return &APIError{Status: http.StatusConflict, Message: msg}
}

func NewUnauthorized(msg string) *APIError {
	return &APIError{Status: http.StatusUnauthorized, Message: msg}
}

func NewInternal(msg string) *APIError {
	return &APIError{Status: http.StatusInternalServerError, Message: msg}
}
