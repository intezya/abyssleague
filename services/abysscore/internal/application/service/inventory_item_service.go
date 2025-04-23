package applicationservice

import (
	"abysscore/internal/domain/dto"
	repositoryports "abysscore/internal/domain/repository"
	"context"
)

type InventoryItemService struct {
	repository     repositoryports.InventoryItemRepository
	userRepository repositoryports.UserRepository
}

func NewInventoryItemService(
	repository repositoryports.InventoryItemRepository,
	userRepository repositoryports.UserRepository,
) *InventoryItemService {
	return &InventoryItemService{repository: repository, userRepository: userRepository}
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
	_, err := i.repository.FindByUserIDAndID(ctx, userID, inventoryItemID)

	if err != nil {
		return err
	}

	err = i.repository.Delete(ctx, inventoryItemID)

	if err != nil {
		return err
	}

	// TODO: send notification to user

	return nil
}

func (i *InventoryItemService) SetInventoryItemAsCurrent(
	ctx context.Context,
	user *dto.UserDTO,
	inventoryItemID int,
) error {
	item, err := i.repository.FindByUserIDAndID(ctx, user.ID, inventoryItemID)

	if err != nil {
		return err // item not found in inventory
	}

	err = i.userRepository.SetInventoryItemAsCurrent(ctx, user, item)

	if err != nil {
		return err // ???
	}

	return nil
}
