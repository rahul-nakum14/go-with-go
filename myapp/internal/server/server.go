package server

import (
	"log"
	"net/http"
	"myapp/internal/routes"
	"myapp/internal/middleware"
)

func NewServer(addr string) *http.Server {
	mux := http.NewServeMux()

	routes.RegisterRoutes(mux)

	handler := middleware.Logger(mux)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

func Start(addr string) {
	srv := NewServer(addr)
	log.Printf("Server running on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
