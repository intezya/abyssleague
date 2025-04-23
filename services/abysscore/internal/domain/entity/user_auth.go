package entity

import (
	"time"
)

type passwordComparator func(raw, hash string) bool

// AuthenticationData represents user authentication state
type AuthenticationData struct {
	id           int
	username     string
	password     string
	hwid         *string
	blockedUntil *time.Time
	blockReason  *string
}

func NewAuthenticationData(
	id int,
	username string,
	password string,
	hwid *string,
	blockedUntil *time.Time,
	blockReason *string,
) *AuthenticationData {
	return &AuthenticationData{
		id:           id,
		username:     username,
		password:     password,
		hwid:         hwid,
		blockedUntil: blockedUntil,
		blockReason:  blockReason,
	}
}

func (a *AuthenticationData) ComparePassword(password string, comparator passwordComparator) bool {
	return comparator(password, a.password)
}

func (a *AuthenticationData) CompareHardwareID(
	hardwareID string,
	comparator passwordComparator,
) (ok bool, needsUpdate bool) {
	if a.hwid == nil {
		return true, true
	}

	return comparator(hardwareID, *a.hwid), false
}

func (a *AuthenticationData) IsAccountLocked() bool {
	return a.blockedUntil != nil && a.blockedUntil.After(time.Now())
}

func (a *AuthenticationData) TokenData() *TokenData {
	var hwid string

	if a.hwid != nil {
		hwid = *a.hwid
	}

	return &TokenData{
		ID:       a.id,
		Username: a.username,
		Hwid:     hwid,
	}
}

func (a *AuthenticationData) SetHWID(hwid string) {
	a.hwid = &hwid
}

func (a *AuthenticationData) UserID() int {
	return a.id
}

func (a *AuthenticationData) BlockReason() *string {
	return a.blockReason
}
