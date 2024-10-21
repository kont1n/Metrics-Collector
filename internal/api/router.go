package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) InitRoutes() *chi.Mux {
	h.loger.Debugln("InitRoutes start")
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(h.LogAPI)
	router.Use(h.gzipMiddleware)

	router.Route("/update", func(r chi.Router) {
		r.Post("/", h.postJSONMetric)
		r.Post("/{type}/{metric}/{value}", h.postMetric)
	})

	router.Route("/value", func(r chi.Router) {
		r.Post("/", h.getJSONMetric)
		r.Get("/{type}/{metric}", h.getMetric)
	})

	router.Get("/", h.indexHandler)

	h.loger.Debugln("InitRoutes end")
	return router
}
