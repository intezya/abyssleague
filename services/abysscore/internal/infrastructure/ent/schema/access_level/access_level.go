package access_level

import (
	"database/sql/driver"
	"fmt"
	"github.com/intezya/pkglib/logger"
	"strings"
)

//go:generate stringer -type=AccessLevel
type AccessLevel int

const (
	User AccessLevel = iota
	ViewAllUsers
	ViewInventory
	ViewMatches
	Admin
	CreateItem
	GiveItem
	RevokeItem
	UpdateItem
	ResetHwid
	AddAdmin
	DeleteItem
	Dev
)

// Value implements the driver.Valuer interface for saving to DB (as string)
func (a AccessLevel) Value() (driver.Value, error) {
	return a.String(), nil
}

// Scan implements the sql.Scanner interface for reading from DB (as string)
func (a *AccessLevel) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		for i := User; i <= Dev; i++ {
			if strings.EqualFold(i.String(), v) {
				*a = i
				return nil
			}
		}
		logger.Log.Debugf("unknown AccessLevel: %v", value)
		return fmt.Errorf("unknown AccessLevel: %v", value)
	case []byte:
		return a.Scan(string(v))
	case int64:
		if v >= int64(User) && v <= int64(Dev) {
			*a = AccessLevel(v)
			return nil
		}
		logger.Log.Debugf("AccessLevel out of range: %v", value)
		return fmt.Errorf("AccessLevel out of range: %v", value)
	default:
		logger.Log.Debugf("invalid AccessLevel value type: %T", value)
		return fmt.Errorf("invalid AccessLevel value type: %T", value)
	}
}

// FromStringOrDefault converts string to AccessLevel (optional helper)
func FromStringOrDefault(s string) AccessLevel {
	for i := User; i <= Dev; i++ {
		if strings.EqualFold(i.String(), s) {
			return i
		}
	}
	return 0
}
