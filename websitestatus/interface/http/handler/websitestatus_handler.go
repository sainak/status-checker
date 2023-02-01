package handler

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/sainak/status-checker/api"
	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/core/logger"
	"github.com/sainak/status-checker/core/myerrors"
)

func ErrBadRequest(s string) error {
	return &myerrors.Error{
		StatusCode: http.StatusBadRequest,
		Status:     "bad_request",
		Message:    "invalid request: " + s,
	}
}

type WebsiteStatusHandler struct {
	Service domain.WebsiteStatusService
}

type WebsiteRequest struct {
	URL string `json:"url"`
}

func (w WebsiteRequest) Bind(r *http.Request) error {
	if w.URL == "" {
		return ErrBadRequest("url is required")
	}
	return nil
}

func (h *WebsiteStatusHandler) CreateWebsite(w http.ResponseWriter, r *http.Request) {

	data := &WebsiteRequest{}
	if err := render.Bind(r, data); err != nil {
		logger.Error(err)
		if err == io.EOF {
			err = ErrBadRequest("empty request body")
		}
		api.RespondForError(w, r, err)
		return
	}
	website := &domain.Website{
		URL:     data.URL,
		AddedAt: zero.NewTime(time.Now(), true),
	}
	err := h.Service.CreateWebsite(r.Context(), website)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, website)
}

func (h *WebsiteStatusHandler) GetAllSites(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	limit := func() int {
		_limit := r.URL.Query().Get("limit")
		ret, err := strconv.Atoi(_limit)
		if _limit == "" || err != nil {
			return 10
		}
		return ret
	}()
	websites, nextCursor, err := h.Service.ListWebsitesStatus(r.Context(), cursor, int64(limit), nil)
	if err != nil {
		api.RespondForError(w, r, err)
		return
	}
	if nextCursor != "" {
		w.Header().Set("X-Next-Cursor", nextCursor)
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, websites)
}
