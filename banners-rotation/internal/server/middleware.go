package server

import (
	"net/http"
	"time"

	"github.com/ezhk/golang-learning/banners-rotation/internal/logger"
	"go.uber.org/zap"
)

type responseObserver struct {
	http.ResponseWriter
	Status int
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	o.Status = code
}

func ProxyLoggerMiddleware(log *logger.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		o := &responseObserver{ResponseWriter: w}
		h.ServeHTTP(o, r)
		spendTime := time.Since(startTime).Seconds()

		log.Info("HTTP request",
			zap.String("source", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.Int("status_code", o.Status),

			zap.String("user_agent", r.Header.Get("User-Agent")),
			zap.Float64("request_time", spendTime),
		)
	})
}
