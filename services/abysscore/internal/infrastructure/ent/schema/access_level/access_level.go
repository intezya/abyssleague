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
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid AccessLevel value: %v", value)
	}

	*a = FromStringOrDefault(strVal)

	return nil
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
