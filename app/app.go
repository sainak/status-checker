package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sainak/status-checker/core/config"
	"github.com/sainak/status-checker/core/logger"
	"github.com/sainak/status-checker/helpers"
	_rootRouter "github.com/sainak/status-checker/root/interface/http/router"
	"net/http"
)

func Run() {
	config.GetConfig()

	db := sqlx.MustOpen("postgres", config.GetDBurl())
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}(db)

	err := db.Ping()
	if err != nil {
		logger.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.NoCache)

	//r.Mount("/debug", middleware.Profiler())
	_rootRouter.RegisterRoutes(r)

	logger.Print(helpers.GetRegisteredRoutes(r))
	server := &http.Server{
		Addr:    "localhost" + config.GetServerAddress(),
		Handler: r,
	}
	logger.Info("Listening on http://" + server.Addr)
	logger.Fatal(server.ListenAndServe())
}
