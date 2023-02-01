package myerrors

import "errors"

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your Item already exist")
	// ErrBadInputParam will throw if the given request-body or params is not valid
	ErrBadInputParam = errors.New("given Param is not valid")
	// ErrBadCursor will throw if the passed cursor is invalid
	ErrBadCursor = errors.New("given cursor is invalid")
)
