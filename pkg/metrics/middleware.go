package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Tortik3000/service-order/pkg/logger"
	httpMetrics "github.com/Tortik3000/service-order/pkg/metrics/http"
	rateLimitMetrics "github.com/Tortik3000/service-order/pkg/metrics/rate_limit"
)

type Middleware interface {
	Metrics(next http.Handler) http.Handler
}

type middleware struct {
	logs logger.Logger
}

func New(
	logs logger.Logger,
) *middleware {
	return &middleware{
		logs: logs,
	}
}

func (m *middleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs := m.logs.With(
			logger.NewField("method", r.Method),
			logger.NewField("path", r.URL.Path),
		)

		start := time.Now()
		logs.Info("start",
			logger.NewField("time", start),
		)

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(rw.statusCode)

		logs.Info("finish",
			logger.NewField("time", time.Now()),
			logger.NewField("status", rw.statusCode),
			logger.NewField("duration", duration),
		)

		path := basePath(r.URL.Path)
		httpMetrics.HTTPRequestDuration.WithLabelValues(r.Method, path, statusCode).Observe(duration)
		httpMetrics.HTTPRequestTotal.WithLabelValues(r.Method, path, statusCode).Inc()

		if rw.statusCode >= 400 {
			httpMetrics.HTTPErrorTotal.WithLabelValues(r.Method, path, statusCode).Inc()
		}

		if rw.statusCode == http.StatusTooManyRequests {
			rateLimitMetrics.RateLimitTotal.WithLabelValues(r.Method, path).Inc()
		}
	})
}

func basePath(p string) string {
	// Simple basePath implementation, can be expanded if needed
	return p
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
