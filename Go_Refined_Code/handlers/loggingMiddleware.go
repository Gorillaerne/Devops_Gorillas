// Package handlers loggingMiddleware
package handlers

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code written by a handler.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs every API request as a structured JSON event.
// It records method, path, status code, duration, and remote IP.
// Sensitive fields (passwords, tokens) are never included.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newResponseWriter(w)

		next.ServeHTTP(wrapped, r)

		status := wrapped.statusCode
		msg := "API request"

		switch { //nolint:gosec
		case status >= 500:
			slog.Error(msg, //nolint:gosec
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", status),
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
				slog.String("remote_ip", r.RemoteAddr),
			)
		case status >= 400:
			slog.Warn(msg, //nolint:gosec
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", status),
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
				slog.String("remote_ip", r.RemoteAddr),
			)
		default:
			slog.Info(msg, //nolint:gosec
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", status),
				slog.Int64("duration_ms", time.Since(start).Milliseconds()),
				slog.String("remote_ip", r.RemoteAddr),
			)
		}
	})
}
