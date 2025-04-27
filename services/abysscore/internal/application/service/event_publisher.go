package applicationservice

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/event"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	eventlib "github.com/intezya/abyssleague/services/abysscore/pkg/event"
)

const eventConsumerWorkerCount = 2
const eventBufferSize = 100

func NewApplicationEventPublisher(
	mainClientNotificationService domainservice.NotificationService,
	// draftClientNotificationService domainservice.NotificationService,
) *eventlib.ApplicationEventPublisher {
	publisher := eventlib.NewApplicationEventPublisher(
		eventConsumerWorkerCount,
		eventBufferSize,
	)

	// Inventory item events
	inventoryItemEventHandlers := event.NewInventoryItemHandlers(mainClientNotificationService)

	publisher.Register(
		event.InventoryItemObtainedEvent{},
		inventoryItemEventHandlers.InventoryItemObtainedEventHandler,
		//middleware...,
	)

	return publisher
}
