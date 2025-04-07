package jwt

import "github.com/golang-jwt/jwt/v5"

type AuthenticationData interface {
	GetID() (id int)
	GetUsername() (username string)
	GetHardwareID() (hwid string)
}

type Claim struct {
	Username   string `json:"username"`
	UserID     int    `json:"user_id"`
	HardwareID string `json:"hardware_id"`
	jwt.RegisteredClaims
}

func (c Claim) GetID() (id int) {
	return c.UserID
}

func (c Claim) GetUsername() (username string) {
	return c.Username
}

func (c Claim) GetHardwareID() (hwid string) {
	return c.HardwareID
}
