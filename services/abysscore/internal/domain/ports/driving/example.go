package drivingports

/*
Интерфейсы, через которые внешние системы могут взаимодействовать с вашим приложением
Обычно содержат методы, представляющие доступную функциональность вашего приложения
*/

// Change any type to entity type

type UserService interface {
	GetUser(id string) (any, error)
	CreateUser(user any) error
	UpdateUser(user any) error
}
