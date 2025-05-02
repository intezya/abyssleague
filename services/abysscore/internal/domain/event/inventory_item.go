package event

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/websocketmessage"
	eventlib "github.com/intezya/abyssleague/services/abysscore/pkg/event"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
	"github.com/intezya/pkglib/logger"
)

type InventoryItemObtainedEvent struct {
	*BaseEvent
	// metadata
	performer optional.Optional[*dto.UserDTO]
	userID    int
	item      *dto.InventoryItemDTO
}

func NewInventoryItemObtainedEvent(
	eventID optional.String,
	performer optional.Optional[*dto.UserDTO],
	userID int,
	item *dto.InventoryItemDTO,
) *InventoryItemObtainedEvent {
	senderID := eventlib.SystemIsSender

	if performer.IsSet() {
		senderID = performer.MustValue().ID
	}

	return &InventoryItemObtainedEvent{
		BaseEvent: newBaseEvent(eventID, userID, senderID),
		performer: performer,
		userID:    userID,
		item:      item,
	}
}

type InventoryItemHandlers struct {
	notificationService domainservice.NotificationService
}

func NewInventoryItemHandlers(
	notificationService domainservice.NotificationService,
) *InventoryItemHandlers {
	return &InventoryItemHandlers{notificationService: notificationService}
}

func (h *InventoryItemHandlers) InventoryItemObtainedEventHandler(event eventlib.ApplicationEvent) {
	typedEvent, ok := event.(*InventoryItemObtainedEvent)
	if !ok {
		return
	}

	performerName := websocketmessage.SystemIsSenderName

	if typedEvent.performer.IsSet() {
		performerName = typedEvent.performer.MustValue().Username
	}

	message := websocketmessage.NewInventoryItemObtainedMessage(
		typedEvent.id,
		performerName,
		typedEvent.item,
	)

	err := h.notificationService.SendToUser(typedEvent.receiverID, message)
	if err != nil {
		logger.Log.Warnln("failed to send message to user:", err)
	}
}
