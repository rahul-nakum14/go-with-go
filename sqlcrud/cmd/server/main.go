package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahulnakum/sqlcrud/internal/config"
	"github.com/rahulnakum/sqlcrud/internal/db"
	"github.com/rahulnakum/sqlcrud/internal/handler"
	"github.com/rahulnakum/sqlcrud/internal/repository"
	"github.com/rahulnakum/sqlcrud/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	appmw "github.com/rahulnakum/sqlcrud/internal/middleware"

)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := repository.NewUserRepo(database, cfg.DBDriver)
	service := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(service)

	r := chi.NewRouter()

	r.Use(appmw.RequestLogger)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Get("/{id}", userHandler.Get)
		r.Put("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)

		// This is how we group and make protected and unprotected routes
		r.Group(func(r chi.Router) {
			r.Put("/{id}", userHandler.Update)
			r.Delete("/{id}", userHandler.Delete)
		})
	})
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
