package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Avik2024/erebus/backend/internal/config"
	"github.com/Avik2024/erebus/backend/internal/health"
	"github.com/Avik2024/erebus/backend/internal/logging"
	"github.com/Avik2024/erebus/backend/internal/metrics"
	"github.com/Avik2024/erebus/backend/internal/version"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	// ----------------------------
	// Load configuration
	// ----------------------------
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ----------------------------
	// Create structured logger
	// ----------------------------
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Inject logger into internal packages
	health.SetLogger(logger)
	version.SetLogger(logger)

	// ----------------------------
	// Initialize build info metric
	// ----------------------------
	metrics.InitBuildInfo(
		version.GetVersion(),
		version.GetCommit(),
		version.GetDate(),
	)

	// ----------------------------
	// Create router & middlewares
	// ----------------------------
	r := chi.NewRouter()
	r.Use(middleware.RequestID)             // generate request ID
	r.Use(middleware.RealIP)                // get real client IP
	r.Use(middleware.Recoverer)             // recover from panics
	r.Use(logging.LoggerMiddleware(logger)) // structured logging
	r.Use(metrics.InstrumentHandler)        // Prometheus metrics with request_id

	// ----------------------------
	// API Endpoints
	// ----------------------------
	r.Get("/api/healthz", health.Handler)
	r.Get("/api/version", version.Handler)

	// Metrics endpoint for Prometheus
	metrics.RegisterMetricsEndpoint(r)

	// Root endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Erebus"))
	})

	// ----------------------------
	// Create HTTP server
	// ----------------------------
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port, // <-- use cfg.App.Port
		Handler: r,
	}

	// ----------------------------
	// Start server
	// ----------------------------
	go func() {
		logger.Info("starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen failed", zap.Error(err))
		}
	}()

	// ----------------------------
	// Graceful shutdown
	// ----------------------------
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}
