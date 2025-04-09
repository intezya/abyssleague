package middleware

import (
	"abysslib/logger"
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"
	"websocket/internal/adapters/controller/http/routes"
	"websocket/internal/infrastructure/metrics"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Debugf("Context keys: %+v", r.Context())
			logger.Log.Debugf("Request ID from context: '%s'", GetRequestID(r.Context()))

			start := time.Now()
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)

			requestDuration := time.Since(start)

			if !routes.IsInfoLogging(r.URL.Path) {
				logger.Log.Debugw(
					"http request",
					"request_id", GetRequestID(r.Context()),
					"method", r.Method,
					"url", r.URL.String(),
					"status", lrw.statusCode,
					"duration", requestDuration,
					"remote_addr", r.RemoteAddr,
					"request_uri", r.RequestURI,
					"body", r.Body, // TODO: body is always empty
					"host", r.Host,
					"method", r.Method,
					"referer", r.Referer(), // TODO: referer is always empty
					"user_agent", r.UserAgent(),
					"content_length", r.ContentLength,
				)
				return
			}

			metrics.ApiRequestDuration.WithLabelValues(
				r.Method,
				r.URL.Path,
			).Observe(float64(requestDuration.Milliseconds()))

			logger.Log.Infow(
				"http request",
				"request_id", GetRequestID(r.Context()),
				"method", r.Method,
				"url", r.URL.String(),
				"status", lrw.statusCode,
				"duration", requestDuration,
				"remote_addr", r.RemoteAddr,
				"request_uri", r.RequestURI,
				"body", r.Body, // TODO: body is always empty
				"host", r.Host,
				"method", r.Method,
				"referer", r.Referer(), // TODO: referer is always empty
				"user_agent", r.UserAgent(),
				"content_length", r.ContentLength,
			)
		},
	)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := lrw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("underlying ResponseWriter does not support Hijack")
}
