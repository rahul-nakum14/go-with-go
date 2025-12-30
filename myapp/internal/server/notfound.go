package server

import (
	"net/http"
)


func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 page not found"))
}


func WithNotFound(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)

		if pattern == "/" && r.URL.Path != "/" {
			NotFoundHandler(w, r)
			return
		}

		mux.ServeHTTP(w, r)
	})
}
