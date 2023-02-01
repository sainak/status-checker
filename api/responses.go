package api

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sainak/status-checker/core/myerrors"
)

// ErrorResponse represent the response error struct
type ErrorResponse struct {
	Message string `json:"message"`
}

func RespondForError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, myerrors.GetErrStatusCode(err))
	render.JSON(w, r, ErrorResponse{err.Error()})
}
