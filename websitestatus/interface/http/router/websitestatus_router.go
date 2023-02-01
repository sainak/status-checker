package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sainak/status-checker/core/domain"
	"github.com/sainak/status-checker/websitestatus/interface/http/handler"
)

func RegisterRoutes(router *chi.Mux, service domain.WebsiteStatusService) {
	h := &handler.WebsiteStatusHandler{Service: service}

	router.Get("/sites", h.GetAllSites)
	router.Post("/sites", h.CreateWebsite)
}
