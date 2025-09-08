package metrics

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// -------------------
// Metrics Definitions
// -------------------

var (
	// Request metrics
	RequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "erebus_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status", "request_id"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "erebus_http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method", "status", "request_id"},
	)

	// Error counters
	Error4xx = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "erebus_http_errors_4xx_total",
			Help: "Total number of client errors (4xx)",
		},
		[]string{"path", "method", "request_id"},
	)

	Error5xx = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "erebus_http_errors_5xx_total",
			Help: "Total number of server errors (5xx)",
		},
		[]string{"path", "method", "request_id"},
	)

	// Build info metric
	BuildInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "erebus_build_info",
			Help: "Erebus build info: 1 for current build",
		},
		[]string{"commit", "version", "date"},
	)
)

// -------------------
// Init Build Info
// -------------------

// InitBuildInfo sets build info metric
func InitBuildInfo(commit, version, date string) {
	BuildInfo.WithLabelValues(commit, version, date).Set(1)
}

// -------------------
// Middleware
// -------------------

// InstrumentHandler wraps http.Handler to record metrics
func InstrumentHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, status: 200}
		requestID := middleware.GetReqID(r.Context())

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			RequestDuration.WithLabelValues(r.URL.Path, r.Method, http.StatusText(recorder.status), requestID).Observe(v)
		}))
		defer timer.ObserveDuration()

		next.ServeHTTP(recorder, r)

		RequestCount.WithLabelValues(r.URL.Path, r.Method, http.StatusText(recorder.status), requestID).Inc()

		// Increment error counters if necessary
		if recorder.status >= 400 && recorder.status < 500 {
			Error4xx.WithLabelValues(r.URL.Path, r.Method, requestID).Inc()
		} else if recorder.status >= 500 {
			Error5xx.WithLabelValues(r.URL.Path, r.Method, requestID).Inc()
		}
	})
}

// -------------------
// Helpers
// -------------------

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// RegisterMetricsEndpoint exposes /metrics
func RegisterMetricsEndpoint(r chi.Router) {
	r.Handle("/metrics", promhttp.Handler())
}
