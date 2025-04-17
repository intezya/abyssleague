package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Hwid     string `json:"hwid"`
}

type JWTConfiguration struct {
	secretKey      []byte
	issuer         string
	expirationTime time.Duration
}

func NewJWTConfiguration(secretKey string, issuer string, expirationTime time.Duration) *JWTConfiguration {
	return &JWTConfiguration{secretKey: []byte(secretKey), issuer: issuer, expirationTime: expirationTime}
}

type Claim struct {
	AuthenticationData *TokenData `json:"authentication_data"`
	jwt.RegisteredClaims
}

type JWTHelper struct {
	*JWTConfiguration
}

func NewJWTHelper(configuration *JWTConfiguration) *JWTHelper {
	return &JWTHelper{JWTConfiguration: configuration}
}

func (j *JWTHelper) TokenGenerator(tokenData *TokenData) string {
	claims := &Claim{
		AuthenticationData: tokenData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(j.secretKey)

	return tokenString
}

func (j *JWTHelper) ValidateToken(tokenString string) (*TokenData, error) {
	claims := &Claim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return j.secretKey, nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithStrictDecoding(),
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims.AuthenticationData, nil
}
