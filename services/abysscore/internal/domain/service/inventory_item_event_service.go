package domainservice

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
)

type InventoryItemEventService interface {
	HandleItemObtained(
		ctx context.Context,
		receiverID int,
		performer optional.Optional[*dto.UserDTO],
		item *dto.InventoryItemDTO,
	)
	HandleItemRevoked(
		ctx context.Context,
		receiverID int,
		performer optional.Optional[*dto.UserDTO],
		item *dto.InventoryItemDTO,
	)
}
