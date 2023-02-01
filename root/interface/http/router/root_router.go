package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sainak/status-checker/root/interface/http/handler"
)

func RegisterRoutes(router *chi.Mux) {
	router.Get("/", handler.Root)
	router.Get("/ping", handler.Ping)
	router.Get("/panic", handler.Panic)
}
