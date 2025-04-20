package repositoryports

import (
	"abysscore/internal/domain/dto"
	"context"
)

type InventoryItemRepository interface {
	Create(ctx context.Context, inventoryItem *dto.CreateInventoryItemDTO) (*dto.InventoryItemDTO, error)
	ExistsByUserIDAndID(ctx context.Context, userID, id int) bool
	FindByUserID(ctx context.Context, userID int) ([]*dto.InventoryItemDTO, error) // TODO: maybe add pagination
	Delete(ctx context.Context, inventoryItemID int) error
}
