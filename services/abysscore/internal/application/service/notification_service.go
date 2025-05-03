package applicationservice

import (
	"context"
	"encoding/json"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/clients"
	"github.com/intezya/pkglib/logger"
)

type NotificationService struct {
	websocketService clients.WebsocketMessagingClient
}

func NewNotificationService(
	websocketService clients.WebsocketMessagingClient,
) *NotificationService {
	return &NotificationService{websocketService: websocketService}
}

func (n NotificationService) SendToUser(userID int, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		logger.Log.Warn("failed to marshal message: %v", err)

		return err
	}

	err = n.websocketService.SendMessage(context.Background(), userID, payload)
	if err != nil {
		logger.Log.Debug("failed to send message: %v", err)

		return err
	}

	return nil
}
