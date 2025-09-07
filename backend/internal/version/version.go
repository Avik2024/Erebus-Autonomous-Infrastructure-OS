package version

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var logger *zap.Logger

// Version value (bump on releases)
var Version = "v0.0.1"

// SetLogger allows main to inject a zap.Logger
func SetLogger(l *zap.Logger) { logger = l }

// Handler returns the current version (with newline)
func Handler(w http.ResponseWriter, r *http.Request) {
	if logger != nil {
		logger.Info("version request", zap.String("version", Version))
	}
	fmt.Fprintln(w, Version)
}