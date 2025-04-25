package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/intezya/pkglib/logger"
	"net/http"
)

type key string

const requestIDKey key = "request_id"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Log.Debug("Starting")

			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
				logger.Log.Debugf("Generated new ID: %s", requestID)
			} else {
				logger.Log.Debugf("Using existing ID: %s", requestID)
			}

			w.Header().Set("X-Request-ID", requestID)

			ctx := context.WithValue(r.Context(), requestIDKey, requestID)

			logger.Log.Debug("Calling next handler")
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
