package ws

import (
	"abysslib/jwt"
	"abysslib/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type SecurityMiddleware struct {
	jwtService jwt.Validate
}

func NewMiddleware(jwtService jwt.Validate) *SecurityMiddleware {
	return &SecurityMiddleware{jwtService: jwtService}
}

func (m *SecurityMiddleware) JwtAuth(w http.ResponseWriter, r *http.Request) (authData jwt.AuthenticationData) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	token := jwt.ExtractFromHeader(authHeader)
	authData, err := m.jwtService.Validate(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	return authData
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)

			logger.Log.Info(
				"http request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()),
				zap.Int("status", lrw.statusCode),
				zap.Duration("duration", time.Since(start)),
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
