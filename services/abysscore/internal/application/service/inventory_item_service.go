package applicationservice

import (
	applicationerror "abysscore/internal/common/errors/application"
	"abysscore/internal/domain/dto"
	repositoryports "abysscore/internal/domain/repository"
	"context"
)

type InventoryItemService struct {
	repository repositoryports.InventoryItemRepository
}

func NewInventoryItemService(repository repositoryports.InventoryItemRepository) *InventoryItemService {
	return &InventoryItemService{repository: repository}
}

func (i *InventoryItemService) GrantToUserByAdmin(
	ctx context.Context,
	userID int,
	itemID int,
	performer *dto.UserDTO,
) (*dto.InventoryItemDTO, error) {
	createRequest := &dto.CreateInventoryItemDTO{
		UserID:         userID,
		ItemID:         itemID,
		ReceivedFromID: performer.ID,
	}

	result, err := i.repository.Create(ctx, createRequest)

	if err != nil {
		return nil, err
	}

	// TODO: send notification to user

	return result, nil
}

func (i *InventoryItemService) FindAllByUserID(ctx context.Context, userID int) ([]*dto.InventoryItemDTO, error) {
	result, err := i.repository.FindByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (i *InventoryItemService) RevokeByAdmin(
	ctx context.Context,
	userID int,
	inventoryItemID int,
	performer *dto.UserDTO,
) error {
	if !i.repository.ExistsByUserIDAndID(ctx, userID, inventoryItemID) {
		return applicationerror.ErrItemNotFoundInInventory
	}

	err := i.repository.Delete(ctx, inventoryItemID)

	if err != nil {
		return err
	}

	// TODO: send notification to user

	return nil
}
