package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Service struct {
	secretKey      []byte
	issuer         string
	expirationTime time.Duration
}

func New(config Configuration) *Service {
	return &Service{
		secretKey:      config.SecretKey(),
		issuer:         config.Issuer(),
		expirationTime: config.ExpirationTime(),
	}
}

func (s Service) Validate(tokenString string) (authData AuthenticationData, err error) {
	claims := &Claim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return s.secretKey, nil
		},
		jwt.WithIssuer(s.issuer),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s Service) Generate(authData AuthenticationData) (tokenString string, err error) {
	claims := &Claim{
		UserID:     authData.GetID(),
		Username:   authData.GetUsername(),
		HardwareID: authData.GetHardwareID(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secretKey)
}
