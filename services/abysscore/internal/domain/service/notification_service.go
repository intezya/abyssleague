package domainservice

import "context"

type NotificationService interface {
	SendToUser(ctx context.Context, userID int, message interface{}) error
}
