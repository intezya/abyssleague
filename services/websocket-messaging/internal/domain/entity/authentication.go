package entity

import (
	"github.com/intezya/abyssleague/services/websocket-messaging/internal/pkg/auth"
	"strconv"
)

type AuthenticationData struct {
	id         int
	username   string
	hardwareID string
}

func NewAuthenticationData(id int, username string, hardwareID string) *AuthenticationData {
	return &AuthenticationData{
		id:         id,
		username:   username,
		hardwareID: hardwareID,
	}
}

func (a *AuthenticationData) ID() int {
	return a.id
}

func (a *AuthenticationData) Username() string {
	return a.username
}

func (a *AuthenticationData) HardwareID() string {
	return a.hardwareID
}

func (a *AuthenticationData) Encode() map[string]string {
	return map[string]string{
		"id":         strconv.Itoa(a.id),
		"username":   a.username,
		"hardwareID": a.hardwareID,
	}
}

func DecodeToAuthenticationData(data *auth.TokenData) *AuthenticationData {
	return &AuthenticationData{
		id:         data.ID,
		username:   data.Username,
		hardwareID: data.Hwid,
	}
}
