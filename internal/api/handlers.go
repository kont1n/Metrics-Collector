package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"Metrics-Collector/internal/storage"
)

func PostMetric(store *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "metric")
		metricValue := chi.URLParam(r, "value")

		fmt.Println("handler:", metricType, metricName, metricValue)

		if metricName == "" {
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "incorrect metric name", http.StatusNotFound)
			return
		}

		if metricValue == "" {
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		}

		switch metricType {
		case "gauge":
			{
				value, err := strconv.ParseFloat(metricValue, 64)
				if err != nil {
					w.Header().Set("Content-Type", "text/plain")
					http.Error(w, "incorrect metric value", http.StatusBadRequest)
					return
				}
				store.SetGauge(metricName, value)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusCreated)
			}
		case "counter":
			{
				value, err := strconv.ParseInt(metricValue, 10, 64)
				if err != nil {
					w.Header().Set("Content-Type", "text/plain")
					http.Error(w, "incorrect metric value", http.StatusBadRequest)
					return
				}
				store.SetCounter(metricName, value)
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusCreated)
			}
		default:
			{
				w.Header().Set("Content-Type", "text/plain")
				http.Error(w, "incorrect metric type", http.StatusBadRequest)
			}
		}
	}
}
