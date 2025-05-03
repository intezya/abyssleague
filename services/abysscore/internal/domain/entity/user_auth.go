package entity

import (
	"time"

	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
)

type passwordComparator func(raw, hash string) bool

// AuthenticationData represents user authentication state.
type AuthenticationData struct {
	id           int
	username     string
	password     string
	hardwareID   *string
	blockedUntil *time.Time
	blockReason  *string
}

func NewAuthenticationData(
	id int,
	username string,
	password string,
	hardwareID *string,
	blockedUntil *time.Time,
	blockReason *string,
) *AuthenticationData {
	return &AuthenticationData{
		id:           id,
		username:     username,
		password:     password,
		hardwareID:   hardwareID,
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
	if a.hardwareID == nil {
		return true, true
	}

	return comparator(hardwareID, *a.hardwareID), false
}

func (a *AuthenticationData) IsAccountLocked() bool {
	return a.blockedUntil != nil && a.blockedUntil.After(time.Now())
}

func (a *AuthenticationData) TokenData() *TokenData {
	var hwid string

	if a.hardwareID != nil {
		hwid = *a.hardwareID
	}

	return &TokenData{
		ID:       a.id,
		Username: a.username,
		Hwid:     hwid,
	}
}

func (a *AuthenticationData) SetHardwareID(hwid string) {
	a.hardwareID = &hwid
}

func (a *AuthenticationData) UserID() int {
	return a.id
}

func (a *AuthenticationData) BlockReason() optional.String {
	return optional.NewP(a.blockReason)
}
