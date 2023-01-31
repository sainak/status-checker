package helpers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sainak/status-checker/core/logger"
)

func GetRegisteredRoutes(r *chi.Mux) string {
	var routes []string
	err := chi.Walk(r,
		func(
			method string,
			route string,
			handler http.Handler,
			middlewares ...func(http.Handler) http.Handler) error {

			routes = append(routes, fmt.Sprintf("%s\t%s", method, route))
			return nil
		})
	if err != nil {
		logger.Error(err)
	}
	return "Registered routes:\n\t\t" + strings.Join(routes, "\n\t\t")
}
