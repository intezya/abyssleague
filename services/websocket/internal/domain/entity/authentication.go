package entity

type AuthenticationData interface {
	GetID() (id int)
	GetUsername() (username string)
	GetHardwareID() (hwid string)
}
