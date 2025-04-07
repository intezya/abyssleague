package ws

import (
	"abysslib/jwt"
	"net/http"
)

type Middleware struct {
	jwtService jwt.Validate
}

func NewMiddleware(jwtService jwt.Validate) *Middleware {
	return &Middleware{jwtService: jwtService}
}

func (m *Middleware) JwtAuth(w http.ResponseWriter, r *http.Request) (authData jwt.AuthenticationData) {
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
