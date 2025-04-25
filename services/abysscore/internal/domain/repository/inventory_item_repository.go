package repositoryports

import (
	"context"

	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
)

type InventoryItemRepository interface {
	Create(
		ctx context.Context,
		inventoryItem *dto.CreateInventoryItemDTO,
	) (*dto.InventoryItemDTO, error)
	FindByUserIDAndID(ctx context.Context, userID, id int) (*dto.InventoryItemDTO, error)
	FindByUserID(
		ctx context.Context,
		userID int,
	) ([]*dto.InventoryItemDTO, error) // TODO: maybe add pagination
	Delete(ctx context.Context, inventoryItemID int) error
}
