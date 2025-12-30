package routes

import (
	"net/http"
	"myapp/internal/handlers"
	"myapp/internal/middleware"
)

func Chain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/", http.HandlerFunc(handlers.HomeHandler))

	mux.Handle(
		"/demotest",
		Chain(
			http.HandlerFunc(handlers.DemoHandler),
			middleware.Auth,
		),
	)
}
