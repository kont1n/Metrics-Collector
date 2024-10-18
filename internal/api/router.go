package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *APIHandler) InitRoutes() *chi.Mux {
	h.loger.Debugln("InitRoutes")
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(h.LogAPI)

	router.Route("/update", func(r chi.Router) {
		r.Post("/", h.postJSONMetric)
		r.Post("/{type}/{metric}/{value}", h.postMetric)
	})

	router.Get("/value/{type}/{metric}", h.getMetric)
	router.Get("/", h.indexHandler)

	return router
}
