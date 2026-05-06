package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"iam-platform/internal/utils"
)

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
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// ✅ THIS is the correct signature for chi.Use()
func ZapLogger(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			ww := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			reqID := middleware.GetReqID(r.Context())
			userID := utils.GetUserID(r.Context())

			next.ServeHTTP(ww, r)

			log.Info("http request",
				zap.String("request_id", reqID),
				zap.String("user_id", userID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.status),
				zap.Int("bytes", ww.size),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}
