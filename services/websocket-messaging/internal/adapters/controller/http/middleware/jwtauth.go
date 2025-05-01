package middleware

import (
	"net/http"
	"strings"

	"github.com/intezya/abyssleague/services/websocket-messaging/internal/domain/entity"
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/pkg/auth"
	"github.com/intezya/pkglib/logger"
)

type SecurityMiddleware struct {
	jwtService *auth.JWTHelper
}

func NewMiddleware(jwtService *auth.JWTHelper) *SecurityMiddleware {
	return &SecurityMiddleware{jwtService: jwtService}
}

func (m *SecurityMiddleware) JwtAuth(
	w http.ResponseWriter,
	r *http.Request,
) (authenticationData *entity.AuthenticationData) {
	authHeader := r.Header.Get("Authorization")
	var token string

	if authHeader != "" {
		token = softExtractTokenFromHeader(authHeader, "Bearer ", "Token ")
	} else {
		token = r.URL.Query().Get("token")
		if token == "" {
			logger.Log.Debug("missing authorization header and query token")
			http.Error(w, "missing token", http.StatusUnauthorized)
			return nil
		}
	}

	tokenData, err := m.jwtService.ValidateToken(token)
	if err != nil {
		logger.Log.Debug("validating token error: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)

		return authenticationData
	}

	authenticationData = entity.DecodeToAuthenticationData(tokenData)

	if authenticationData == nil {
		logger.Log.Debug("malformed token data")
		http.Error(w, "malformed token data", http.StatusUnauthorized)

		return authenticationData
	}

	return authenticationData
}

func softExtractTokenFromHeader(tokenString string, availablePrefixes ...string) string {
	for _, prefix := range availablePrefixes {
		if token, found := strings.CutPrefix(tokenString, prefix); found {
			return token
		}
	}

	return tokenString
}
