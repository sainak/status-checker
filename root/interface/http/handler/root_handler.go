package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

func Root(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "root")
}

func Ping(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{
		"message":      "pong",
		"current_time": time.Now(),
	})
}

func Panic(_ http.ResponseWriter, _ *http.Request) {
	panic("test")
}
