// Package middleware provides composable HTTP middleware.
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain composes multiple middlewares into one, applied left→right.
// The first middleware is the outermost (runs first on request, last on response).
func Chain(ms ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(ms) - 1; i >= 0; i-- {
			next = ms[i](next)
		}
		return next
	}
}

// ── Request ID ────────────────────────────────────────────────────────────────

type ctxKey string

const RequestIDKey ctxKey = "request_id"

// RequestID injects a unique request ID into every request context and response header.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = uuid.NewString()
			}
			ctx := context.WithValue(r.Context(), RequestIDKey, id)
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ── Structured Logger ─────────────────────────────────────────────────────────

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Logger logs each request with method, path, status, duration, and request ID.
func Logger(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rw, r)

			reqID, _ := r.Context().Value(RequestIDKey).(string)
			log.InfoContext(r.Context(), "request",
				"method",     r.Method,
				"path",       r.URL.Path,
				"status",     rw.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"bytes",      rw.size,
				"request_id", reqID,
				"ip",         realIP(r),
			)
		})
	}
}

// ── Panic Recovery ────────────────────────────────────────────────────────────

// Recover catches any panics in downstream handlers and returns a 500.
func Recover(log *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.ErrorContext(r.Context(), "panic recovered",
						"panic", rec,
						"stack", string(debug.Stack()),
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// ── CORS ──────────────────────────────────────────────────────────────────────

// CORS handles cross-origin requests for the given allowed origins.
func CORS(allowedOrigins []string) Middleware {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if allowed[origin] || allowed["*"] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ── Token-bucket Rate Limiter ─────────────────────────────────────────────────

type bucket struct {
	tokens   float64
	lastSeen time.Time
	mu       sync.Mutex
}

// RateLimit limits each IP to rps requests per second using a token bucket.
func RateLimit(rps int) Middleware {
	clients := sync.Map{}
	rate := float64(rps)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := realIP(r)
			v, _ := clients.LoadOrStore(ip, &bucket{tokens: rate, lastSeen: time.Now()})
			b := v.(*bucket)

			b.mu.Lock()
			now := time.Now()
			elapsed := now.Sub(b.lastSeen).Seconds()
			b.tokens = min(rate, b.tokens+elapsed*rate)
			b.lastSeen = now

			if b.tokens < 1 {
				b.mu.Unlock()
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			b.tokens--
			b.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.SplitN(ip, ",", 2)[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
