package middleware

import (
	"abysslib/logger"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Debugf("Context keys: %+v", r.Context())
			logger.Log.Debugf("Request ID from context: '%s'", GetRequestID(r.Context()))

			start := time.Now()
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)

			logger.Log.Infow(
				"http request",
				"request_id", GetRequestID(r.Context()),
				"method", r.Method,
				"url", r.URL.String(),
				"status", lrw.statusCode,
				"duration", time.Since(start),
				"remoteaddr", r.RemoteAddr,
				"requesturi", r.RequestURI,
				"body", r.Body,
				"host", r.Host,
				"method", r.Method,
				"url", r.URL,
				"referer", r.Referer(),
				"useragent", r.UserAgent(),
				"contentlength", r.ContentLength,
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
