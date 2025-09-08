package version

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

var logger *zap.Logger

// SetLogger allows main to inject a zap.Logger
func SetLogger(l *zap.Logger) { logger = l }

// Version information
var Version = "v0.0.1"
var Commit = "dev"
var Date = "unknown"

// Getters for ldflags injection
func GetVersion() string { return Version }
func GetCommit() string  { return Commit }
func GetDate() string    { return Date }

// Handler returns version info as JSON
func Handler(w http.ResponseWriter, r *http.Request) {
	reqID := middleware.GetReqID(r.Context()) // Correct way to get request ID

	info := map[string]string{
		"version": Version,
		"commit":  Commit,
		"date":    Date,
		"req_id":  reqID,
	}

	if logger != nil {
		logger.Info("version request",
			zap.String("version", Version),
			zap.String("commit", Commit),
			zap.String("date", Date),
			zap.String("req_id", reqID),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode version: %v", err), http.StatusInternalServerError)
	}
}
