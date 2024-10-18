package api

import (
	"fmt"
	"html"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"Metrics-Collector/internal/service"
)

type ApiHandler struct {
	service *service.Service
	loger   *zap.SugaredLogger
}

func NewHandler(service *service.Service, loger *zap.SugaredLogger) *ApiHandler {
	return &ApiHandler{
		service: service,
		loger:   loger,
	}
}

func (h *ApiHandler) postMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "metric")
	metricValue := chi.URLParam(r, "value")

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

		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		}
		h.service.SetGauge(metricName, value)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

	case "counter":

		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "incorrect metric value", http.StatusBadRequest)
			return
		}
		h.service.SetCounter(metricName, value)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

	default:

		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "incorrect metric type", http.StatusBadRequest)

	}
}

func (h *ApiHandler) getMetrics(w http.ResponseWriter, r *http.Request) {
	var answer string
	w.Header().Set("Content-Type", "text/plain")

	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "metric")

	switch metricType {
	case "gauge":
		{
			value, exists := h.service.GetGauge(metricName)
			if !exists {
				http.Error(w, "unknown metric", http.StatusNotFound)
				return
			}
			answer = strconv.FormatFloat(value, 'f', -1, 64)
		}
	case "counter":
		{
			value, exists := h.service.GetCounter(metricName)
			if !exists {
				http.Error(w, "unknown metric", http.StatusNotFound)
				return
			}
			answer = fmt.Sprintf("%d", value)
		}
	default:
		{
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "incorrect metric type", http.StatusBadRequest)
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(answer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ApiHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	var result string

	for metric, value := range h.service.GetCounters() {
		result += fmt.Sprintf("%s: %d\n", metric, value)
	}

	m := h.service.GetGauges()

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		result += fmt.Sprintf("%s: %.3f\n", k, m[k])
	}

	htmlString := "<html><head><title>Metrics</title></head><body><pre>" +
		html.EscapeString(result) +
		"</pre></body></html>"

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(htmlString)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
