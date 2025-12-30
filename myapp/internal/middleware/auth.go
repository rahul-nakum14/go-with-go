package middleware

import (
	"log"
	"net/http"
	"os"
)

func Auth(next http.Handler) http.Handler {
	var logger = log.New(os.Stdout, "AUTH MIDDLEWARE: ", log.LstdFlags)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Auth middleware invoked",r.Header)
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
