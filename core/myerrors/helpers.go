package myerrors

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lib/pq"
)

func ParseDBError(err error) error {
	switch e := err.(type) {
	case *pq.Error:

		//https://github.com/lib/pq/blob/922c00e176fb3960d912dc2c7f67ea2cf18d27b0/error.go#L78
		switch e.Code {
		case "23502":
			// not-null constraint violation
			return &Error{
				StatusCode: http.StatusConflict,
				Status:     "conflict",
				Message:    fmt.Sprint("some required data was left out:", e.Message),
			}

		case "23505":
			// unique constraint violation
			return &Error{
				StatusCode: http.StatusConflict,
				Status:     "conflict",
				Message:    fmt.Sprint("this record already exists:", e.Message),
			}

		}

	case *strconv.NumError:
		return &Error{
			StatusCode: http.StatusBadRequest,
			Status:     "bad_request",
			Message:    fmt.Sprintf(`"%s" is not a valid number.`, e.Num),
		}

	default:
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		}
	}
	return ErrInternalServerError
}
