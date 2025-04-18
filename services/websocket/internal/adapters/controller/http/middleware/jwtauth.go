package middleware

import (
	"github.com/intezya/pkglib/logger"
	"net/http"
	"strings"
	"websocket/internal/domain/entity"
	"websocket/internal/pkg/auth"
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

	if authHeader == "" {
		logger.Log.Debug("missing authorization header")
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	token := softExtractTokenFromHeader(authHeader, "Bearer ", "Token ")
	tokenData, err := m.jwtService.ValidateToken(token)

	if err != nil {
		logger.Log.Debug("validating token error: ", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	authenticationData = entity.DecodeToAuthenticationData(tokenData)

	if authenticationData == nil {
		logger.Log.Debug("malformed token data")
		http.Error(w, "malformed token data", http.StatusUnauthorized)
		return
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
