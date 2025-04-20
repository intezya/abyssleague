package domainservice

import (
	"abysscore/internal/domain/dto"
	"context"
)

type InventoryItemService interface {
	GrantToUserByAdmin(
		ctx context.Context,
		userID int,
		itemID int,
		performer *dto.UserDTO,
	) (*dto.InventoryItemDTO, error)
	//GrantToUserBySystem(ctx context.Context, request *dto.GrantInventoryItemDTO) error

	FindAllByUserID(ctx context.Context, userID int) ([]*dto.InventoryItemDTO, error)

	RevokeByAdmin(ctx context.Context, userID int, inventoryItemID int, performer *dto.UserDTO) error
	//RevokeBySystem(ctx context.Context, inventoryItemID int, performer *dto.UserDTO) error
}
