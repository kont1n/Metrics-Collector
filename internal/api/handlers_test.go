package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"Metrics-Collector/internal/storage"
)

func TestEmptyMetricName(t *testing.T) {
	r := httptest.NewRequest("POST", "/update/gauge/", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	PostMetric(nil)(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEmptyMetricValue(t *testing.T) {
	r := httptest.NewRequest("POST", "/update/gauge/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("metric", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	PostMetric(nil)(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIncorrectMetricType(t *testing.T) {
	r := httptest.NewRequest("POST", "/update/test/test/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "test")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "1")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	PostMetric(nil)(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIncorrectMetricValue(t *testing.T) {
	r := httptest.NewRequest("POST", "/update/gauge/test/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	PostMetric(nil)(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestIncorrectMetricValue2(t *testing.T) {
	r := httptest.NewRequest("POST", "/update/counter/test/test", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "counter")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	PostMetric(nil)(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostGaugeMetric(t *testing.T) {
	store := storage.NewMemStorage()
	r := httptest.NewRequest("POST", "/update/gauge/test/1.0", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "gauge")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "1.0")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	PostMetric(store)(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1.0, store.GetGauge("test"))
}

func TestPostCounterMetric(t *testing.T) {
	store := storage.NewMemStorage()
	r := httptest.NewRequest("POST", "/update/counter/test/1", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("type", "counter")
	rctx.URLParams.Add("metric", "test")
	rctx.URLParams.Add("value", "1")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	PostMetric(store)(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1), store.GetCounter("test"))
}
