package middleware

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	logger := log.New(os.Stdout, "PANIC: ", log.LstdFlags)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Printf("panic recovered: %v\n%s", err, debug.Stack())

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
