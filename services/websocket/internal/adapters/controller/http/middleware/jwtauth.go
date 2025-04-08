package middleware

import (
	"abysslib/jwt"
	"abysslib/logger"
	"net/http"
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
		logger.Log.Debug("missing authorization header")
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	token := jwt.ExtractFromHeader(authHeader)
	authData, err := m.jwtService.Validate(token)
	if err != nil {
		logger.Log.Debug("validating token error: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	return authData
}
