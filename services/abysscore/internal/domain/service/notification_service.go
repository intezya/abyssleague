package domainservice

type NotificationService interface {
	SendToUser(userID int, message interface{}) error
}
