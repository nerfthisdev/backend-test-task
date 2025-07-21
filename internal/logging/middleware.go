package logging

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			next.ServeHTTP(lrw, r)
			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", lrw.status),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}
