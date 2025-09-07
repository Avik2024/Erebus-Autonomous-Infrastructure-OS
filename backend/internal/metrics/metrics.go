package metrics

import (
	"net/http"

	"github.com/go-chi/chi/v5" // <- needed for chi.Router
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "erebus_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "erebus_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status"},
	)
)

// InstrumentHandler wraps any http.Handler and records metrics
func InstrumentHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, status: 200}

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			RequestDuration.WithLabelValues(r.URL.Path, r.Method, http.StatusText(recorder.status)).Observe(v)
		}))
		defer timer.ObserveDuration()

		next.ServeHTTP(recorder, r)
		RequestCount.WithLabelValues(r.URL.Path, r.Method, http.StatusText(recorder.status)).Inc()
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// RegisterMetricsEndpoint adds /metrics route for Prometheus to scrape
func RegisterMetricsEndpoint(r chi.Router) {
	r.Handle("/metrics", promhttp.Handler())
}
