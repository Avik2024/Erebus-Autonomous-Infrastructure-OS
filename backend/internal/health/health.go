package health

import (
	"net/http"

	"go.uber.org/zap"
)

var logger *zap.Logger

// SetLogger allows main to inject a zap.Logger
func SetLogger(l *zap.Logger) { logger = l }

// Handler responds with a simple alive message
func Handler(w http.ResponseWriter, r *http.Request) {
	if logger != nil {
		logger.Info("health check", zap.String("path", r.URL.Path))
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Erebus is alive 🚀"))
}

