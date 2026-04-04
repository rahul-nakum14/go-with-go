package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourorg/go-user-crud/internal/handler"
	"github.com/yourorg/go-user-crud/internal/middleware"
	"github.com/yourorg/go-user-crud/internal/repository/postgres"
	"github.com/yourorg/go-user-crud/internal/service"
	"github.com/yourorg/go-user-crud/pkg/logger"
	"github.com/yourorg/go-user-crud/pkg/validator"
)

// @title           User CRUD API
// @version         1.0
// @description     Industry-standard Go REST API for user management.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// ── Logger ──────────────────────────────────────────────────────────────
	log := logger.New(os.Getenv("LOG_LEVEL"))

	// ── Config ───────────────────────────────────────────────────────────────
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("failed to load config", "error", err)
	}

	// ── Database ─────────────────────────────────────────────────────────────
	db, err := postgres.NewPool(cfg.DSN)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}
	defer db.Close()

	log.Info("database connection established")

	// ── Dependencies ─────────────────────────────────────────────────────────
	validate := validator.New()
	userRepo := postgres.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, log)
	userHandler := handler.NewUserHandler(userSvc, validate, log)

	// ── Router ────────────────────────────────────────────────────────────────
	mux := http.NewServeMux()

	// Middleware stack (applied manually — no framework needed)
	stack := middleware.Chain(
		middleware.Recover(log),
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.CORS(cfg.AllowedOrigins),
		middleware.RateLimit(cfg.RateLimit),
	)

	// Routes — pure stdlib routing (Go 1.22+ enhanced patterns)
	mux.HandleFunc("GET    /health",            handler.HealthCheck)
	mux.HandleFunc("GET    /api/v1/users",      userHandler.List)
	mux.HandleFunc("POST   /api/v1/users",      userHandler.Create)
	mux.HandleFunc("GET    /api/v1/users/{id}", userHandler.GetByID)
	mux.HandleFunc("PUT    /api/v1/users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", userHandler.Delete)

	// ── Server ────────────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      stack(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── Graceful shutdown ─────────────────────────────────────────────────────
	serverErr := make(chan error, 1)
	go func() {
		log.Info("server starting", "addr", srv.Addr)
		serverErr <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", "error", err)
		}
	case sig := <-quit:
		log.Info("shutdown signal received", "signal", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}
