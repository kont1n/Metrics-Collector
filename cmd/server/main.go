package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"Metrics-Collector/internal/api"
	"Metrics-Collector/internal/storage"
)

const (
	DefaultServer = "localhost:8080"
)

var (
	err         error
	flagRunAddr string
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", DefaultServer, "address and port to run server")
	flag.Parse()
}

func main() {
	parseFlags()

	store := storage.NewMemStorage()

	router := chi.NewRouter()
	router.Post("/update/{type}/{metric}/{value}", api.PostMetric(store))
	router.Get("/value/{type}/{metric}", api.GetMetrics(store))
	router.Get("/", api.IndexHandler(store))

	fmt.Printf("Server started on %s\n", flagRunAddr)
	if err = http.ListenAndServe(flagRunAddr, router); err != nil {
		fmt.Println("Web server error:", err.Error())
		return
	}
}
