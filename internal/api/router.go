package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *ApiHandler) InitRoutes() *chi.Mux {
	h.loger.Debugln("InitRoutes")
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(h.LogAPI)

	router.Post("/update/{type}/{metric}/{value}", h.postMetric)
	router.Get("/value/{type}/{metric}", h.getMetrics)
	router.Get("/", h.indexHandler)

	return router
}
