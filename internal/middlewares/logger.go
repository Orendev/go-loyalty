package middlewares

import (
	"net/http"
	"time"

	"github.com/Orendev/go-loyalty/internal/logger"
	"go.uber.org/zap"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &logger.ResponseData{
			Status: 0,
			Size:   0,
		}

		lw := logger.LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Log.Info("got incoming HTTP request and response",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.Status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.Size),
		)

	})
}
