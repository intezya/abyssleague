package access_level

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
