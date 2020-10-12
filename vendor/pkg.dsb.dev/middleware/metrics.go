package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics is a middleware function that manipulates prometheus metrics for requests
// it handles.
func Metrics() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, h)
	}
}
