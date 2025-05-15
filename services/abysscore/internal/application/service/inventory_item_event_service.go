package applicationservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/websocketmessage"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
	"github.com/intezya/pkglib/logger"
	"golang.org/x/sync/errgroup"
)

type InventoryItemEventService struct {
	notificationService domainservice.NotificationService
}

func NewInventoryItemEventService(notificationService domainservice.NotificationService) *InventoryItemEventService {
	return &InventoryItemEventService{notificationService: notificationService}
}

func (s *InventoryItemEventService) HandleItemObtained(
	ctx context.Context,
	receiverID int,
	performer optional.Optional[*dto.UserDTO],
	item *dto.InventoryItemDTO,
) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationEventService.HandleRegistration")
	defer span.End()

	eventID := uuid.NewString()
	tracer.AddAttribute(ctx, "event_id", eventID)

	performerName := websocketmessage.SystemIsSenderName

	if performer.IsSet() {
		performerName = performer.MustValue().Username
	}

	group, _ := errgroup.WithContext(ctx)

	group.Go(
		func() error {
			message := websocketmessage.NewInventoryItemObtainedMessage(
				eventID,
				performerName,
				item,
			)

			return s.notificationService.SendToUser(ctx, receiverID, message)
		},
	)

	err := group.Wait()

	if err != nil {
		logger.Log.Warnln("failed to send message to user:", err)
	}
}

func (s *InventoryItemEventService) HandleItemRevoked(
	ctx context.Context,
	receiverID int,
	performer optional.Optional[*dto.UserDTO],
	item *dto.InventoryItemDTO,
) {
	ctx, span := tracer.StartSpan(ctx, "AuthenticationEventService.HandleItemRevoked")
	defer span.End()

	eventID := uuid.NewString()
	tracer.AddAttribute(ctx, "event_id", eventID)

	performerName := websocketmessage.SystemIsSenderName

	if performer.IsSet() {
		performerName = performer.MustValue().Username
	}

	group, _ := errgroup.WithContext(ctx)

	group.Go(
		func() error {
			message := websocketmessage.NewInventoryItemRevokedMessage(
				eventID,
				performerName,
				item,
			)

			return s.notificationService.SendToUser(ctx, receiverID, message)
		},
	)

	err := group.Wait()

	if err != nil {
		logger.Log.Warnln("failed to send message to user:", err)
	}
}
