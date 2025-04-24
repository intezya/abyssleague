package access_level

import (
	"database/sql/driver"
	"fmt"
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

// Value implements the driver.Valuer interface for saving to DB (as string).
func (a AccessLevel) Value() (driver.Value, error) {
	return a.String(), nil
}

// Scan implements the sql.Scanner interface for reading from DB (as string).
func (a *AccessLevel) Scan(value interface{}) error {
	switch typedValue := value.(type) {
	case string:
		for i := User; i <= Dev; i++ {
			if strings.EqualFold(i.String(), typedValue) {
				*a = i

				return nil
			}
		}

		return fmt.Errorf("unknown AccessLevel: %typedValue", value)
	case []byte:
		return a.Scan(string(typedValue))
	case int64:
		if typedValue >= int64(User) && typedValue <= int64(Dev) {
			*a = AccessLevel(typedValue)

			return nil
		}

		return fmt.Errorf("AccessLevel out of range: %typedValue", value)
	default:
		return fmt.Errorf("invalid AccessLevel value type: %T", value)
	}
}

// FromStringOrDefault converts string to AccessLevel (optional helper).
func FromStringOrDefault(s string) AccessLevel {
	for i := User; i <= Dev; i++ {
		if strings.EqualFold(i.String(), s) {
			return i
		}
	}

	return 0
}
