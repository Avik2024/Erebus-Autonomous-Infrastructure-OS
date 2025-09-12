package main

import (
	"context"
	"encoding/json"
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

// ----------------------------
// Mock Controllers
// ----------------------------
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body map[string]string
	_ = json.NewDecoder(r.Body).Decode(&body)
	resp := map[string]interface{}{
		"message": "User registered successfully",
		"user":    body,
	}
	json.NewEncoder(w).Encode(resp)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body map[string]string
	_ = json.NewDecoder(r.Body).Decode(&body)
	resp := map[string]interface{}{
		"message": "User logged in successfully",
		"user":    body,
		"cookie":  "mock-auth-cookie",
	}
	json.NewEncoder(w).Encode(resp)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"name":  "Test User",
		"email": "test@example.com",
	}
	json.NewEncoder(w).Encode(resp)
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projects := []map[string]interface{}{
		{"id": 1, "title": "Project One", "description": "First project"},
		{"id": 2, "title": "Project Two", "description": "Second project"},
	}
	json.NewEncoder(w).Encode(projects)
}

func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body map[string]string
	_ = json.NewDecoder(r.Body).Decode(&body)
	resp := map[string]interface{}{
		"message": "Project created successfully",
		"project": body,
	}
	json.NewEncoder(w).Encode(resp)
}

// ----------------------------
// Main
// ----------------------------
func main() {
	// ----------------------------
	// Load configuration
	// ----------------------------
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("App config loaded: Env=%s, Port=%s, DB=%s, Redis Enabled=%v",
		cfg.App.Env, cfg.App.Port, cfg.Database.URL, cfg.Redis.Enabled)

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

	// ----------------------------
	// User & Projects Endpoints
	// ----------------------------
	r.Post("/api/register", Register)
	r.Post("/api/login", Login)
	r.Get("/api/profile", Profile)

	r.Get("/api/projects", ListProjects)
	r.Post("/api/projects", CreateProject)

	// Metrics endpoint for Prometheus
	metrics.RegisterMetricsEndpoint(r)

	// Root endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Erebus"))
	})

	// ----------------------------
	// Create HTTP server
	// ----------------------------
	addr := ":" + cfg.App.Port
	srv := &http.Server{
		Addr:    addr,
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
