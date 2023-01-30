package handler

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/sainak/status-checker/core/domain"
	errors2 "github.com/sainak/status-checker/core/errors"
	"net/http"
	"strconv"
	"time"
)

type WebsiteStatusHandler struct {
	Service domain.WebsiteStatusService
}

type WebsiteRequest struct {
	URL string `json:"url"`
}

func (w WebsiteRequest) Bind(r *http.Request) error {
	if w.URL == "" {
		return ErrInvalidRequest("missing url")
	}
	return nil
}

func ErrInvalidRequest(s string) error {
	return errors.New("invalid request: " + s)
}

func (h *WebsiteStatusHandler) CreateWebsite(w http.ResponseWriter, r *http.Request) {

	data := &WebsiteRequest{}
	if err := render.Bind(r, data); err != nil {
		render.JSON(w, r, errors2.ResponseError{Message: err.Error()})
		return
	}
	website := &domain.Website{
		URL:     data.URL,
		AddedAt: time.Now(),
	}
	err := h.Service.CreateWebsite(r.Context(), website)
	if err != nil {
		render.JSON(w, r, errors2.ResponseError{Message: err.Error()})
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
	websites, nextCursor, err := h.Service.ListWebsites(r.Context(), cursor, int64(limit), nil)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, errors2.ResponseError{Message: err.Error()})
		return
	}
	if nextCursor != "" {
		w.Header().Set("X-Next-Cursor", nextCursor)
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, websites)
}
