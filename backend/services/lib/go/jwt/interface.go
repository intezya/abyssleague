package jwt

import "time"

type Validate interface {
	Validate(tokenString string) (authData AuthenticationData, err error)
}

type Generate interface {
	Generate(authData AuthenticationData) (tokenString string, err error)
}

type Configuration interface {
	SecretKey() []byte
	Issuer() string
	ExpirationTime() time.Duration
}

type TokenService interface {
	Validate
	Generate
}
