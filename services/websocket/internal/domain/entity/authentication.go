package entity

import (
	"strconv"
	"websocket/internal/pkg/auth"
)

type AuthenticationData struct {
	id       int
	username string
	hwid     string
}

func NewAuthenticationData(id int, username string, hwid string) *AuthenticationData {
	return &AuthenticationData{
		id:       id,
		username: username,
		hwid:     hwid,
	}
}

func (a *AuthenticationData) ID() int {
	return a.id
}

func (a *AuthenticationData) Username() string {
	return a.username
}

func (a *AuthenticationData) HardwareID() string {
	return a.hwid
}

func (a *AuthenticationData) Encode() map[string]string {
	return map[string]string{
		"id":       strconv.Itoa(a.id),
		"username": a.username,
		"hwid":     a.hwid,
	}
}

func DecodeToAuthenticationData(data *auth.TokenData) *AuthenticationData {
	return &AuthenticationData{
		id:       data.ID,
		username: data.Username,
		hwid:     data.Hwid,
	}
}
