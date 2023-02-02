package myerrors

import (
	"net/http"
)

// Error is a custom error wrapper with more information
type Error struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func GetErrStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch e := err.(type) {
	case *Error:
		return e.StatusCode
	}
	return http.StatusInternalServerError
}

func GetErrStatus(err error) string {
	if err == nil {
		return "ok"
	}

	switch e := err.(type) {
	case *Error:
		return e.Status
	}
	return "internal_server_error"
}

var (
	ErrInternalServerError = &Error{
		StatusCode: http.StatusInternalServerError,
		Status:     "internal_server_error",
		Message:    "internal server error occurred, we are checking...",
	}
	ErrNotFound = &Error{
		StatusCode: http.StatusNotFound,
		Status:     "not_found",
		Message:    "entity not found",
	}
	ErrConflict = &Error{
		StatusCode: http.StatusConflict,
		Status:     "conflict",
		Message:    "database conflict occurred",
	}
	ErrEntityAlreadyExist = &Error{
		StatusCode: http.StatusConflict,
		Status:     "conflict",
		Message:    "entity already exist",
	}
	ErrBadInputParam = &Error{
		StatusCode: http.StatusBadRequest,
		Status:     "bad_request",
		Message:    "bad input param, check your input params",
	}
	ErrBadCursor = &Error{
		StatusCode: http.StatusBadRequest,
		Status:     "bad_request",
		Message:    "bad cursor, send a valid cursor",
	}
)
