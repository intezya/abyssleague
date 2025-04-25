package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity"
	"time"
)

type JWTConfiguration struct {
	secretKey      []byte
	issuer         string
	expirationTime time.Duration
}

func NewJWTConfiguration(secretKey string, issuer string, expirationTime time.Duration) *JWTConfiguration {
	return &JWTConfiguration{secretKey: []byte(secretKey), issuer: issuer, expirationTime: expirationTime}
}

type Claim struct {
	AuthenticationData *entity.TokenData `json:"authentication_data"`
	jwt.RegisteredClaims
}

type JWTHelper struct {
	*JWTConfiguration
}

func NewJWTHelper(configuration *JWTConfiguration) *JWTHelper {
	return &JWTHelper{JWTConfiguration: configuration}
}

func (j *JWTHelper) TokenGenerator(tokenData *entity.TokenData) string {
	claims := &Claim{
		AuthenticationData: tokenData,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   tokenData.Username,
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expirationTime)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString(j.secretKey)

	return tokenString
}

func (j *JWTHelper) ValidateToken(tokenString string) (*entity.TokenData, error) {
	claims := &Claim{} //nolint:exhaustruct

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
		return nil, errors.New("invalid token")
	}

	return claims.AuthenticationData, nil
}
