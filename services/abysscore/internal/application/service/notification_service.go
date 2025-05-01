package applicationservice

import (
	"context"
	"encoding/json"

	websocketpb "github.com/intezya/abyssleague/proto/websocket"
	drivenports "github.com/intezya/abyssleague/services/abysscore/internal/domain/ports/driven"
	"github.com/intezya/pkglib/logger"
)

type NotificationService struct {
	websocketService drivenports.WebsocketService
}

func NewNotificationService(websocketService drivenports.WebsocketService) *NotificationService {
	return &NotificationService{websocketService: websocketService}
}

func (n NotificationService) SendToUser(userID int, message interface{}) {
	payload, err := json.Marshal(message)
	if err != nil {
		logger.Log.Warn("failed to marshal message: %v", err)

		return
	}

	_ = n.websocketService.SendMessage(
		context.Background(),
		&websocketpb.SendMessageRequest{
			UserId:      int64(userID),
			JsonPayload: payload,
		},
	)
}
