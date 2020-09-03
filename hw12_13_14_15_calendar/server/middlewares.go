package server

import (
	"net/http"
	"time"

	logger "github.com/ezhk/golang-learning/hw12_13_14_15_calendar/logger"
	"go.uber.org/zap"
)

type responseObserver struct {
	http.ResponseWriter
	status int
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	o.status = code
}

func LoggerMiddleware(log *logger.BaseLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		o := &responseObserver{ResponseWriter: w}
		next.ServeHTTP(o, r)
		spendTime := time.Since(startTime).Seconds()

		log.Info("HTTP request",
			zap.String("source", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.Int("status_code", o.status),

			zap.String("user_agent", r.Header.Get("User-Agent")),
			zap.Float64("request_time", spendTime),
		)
	})
}

// * IP клиента;
// * дата и время запроса;
// * метод, path и версия HTTP;
// * код ответа;
// * latency (время обработки запроса, посчитанное, например, с помощью middleware);
// * user agent, если есть.
