package applicationservice

import (
	"context"

	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/event"
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	eventlib "github.com/intezya/abyssleague/services/abysscore/pkg/event"
	"github.com/intezya/abyssleague/services/abysscore/pkg/optional"
)

type InventoryItemService struct {
	inventoryItemRepository repositoryports.InventoryItemRepository
	inventoryRepository     repositoryports.InventoryRepository
	eventPublisher          eventlib.Publisher
}

func NewInventoryItemService(
	repository repositoryports.InventoryItemRepository,
	inventoryRepository repositoryports.InventoryRepository,
	eventPublisher eventlib.Publisher,
) *InventoryItemService {
	return &InventoryItemService{
		inventoryItemRepository: repository,
		inventoryRepository:     inventoryRepository,
		eventPublisher:          eventPublisher,
	}
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

	result, err := tracer.TraceFnWithResult(
		ctx,
		"inventoryItemRepository.Create",
		func(ctx context.Context) (*dto.InventoryItemDTO, error) {
			return i.inventoryItemRepository.Create(ctx, createRequest)
		},
	)
	if err != nil {
		return nil, err
	}

	i.eventPublisher.Publish(event.NewInventoryItemObtainedEvent(
		optional.EmptyOptional[string](),
		optional.New(performer),
		userID,
		result,
	))

	return result, nil
}

func (i *InventoryItemService) FindAllByUserID(
	ctx context.Context,
	userID int,
) ([]*dto.InventoryItemDTO, error) {
	items, err := tracer.TraceFnWithResult(
		ctx,
		"inventoryItemRepository.FindByUserID",
		func(ctx context.Context) ([]*dto.InventoryItemDTO, error) {
			return i.inventoryItemRepository.FindByUserID(ctx, userID)
		},
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (i *InventoryItemService) RevokeByAdmin(
	ctx context.Context,
	userID int,
	inventoryItemID int,
	performer *dto.UserDTO,
) error {
	_, err := tracer.TraceFnWithResult(
		ctx,
		"inventoryItemRepository.FindByUserIDAndID",
		func(ctx context.Context) (*dto.InventoryItemDTO, error) {
			return i.inventoryItemRepository.FindByUserIDAndID(ctx, userID, inventoryItemID)
		},
	)
	if err != nil {
		return err
	}

	err = tracer.TraceFn(
		ctx,
		"inventoryItemRepository.Delete",
		func(ctx context.Context) error {
			return i.inventoryItemRepository.Delete(ctx, inventoryItemID)
		},
	)
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
	item, err := tracer.TraceFnWithResult(
		ctx,
		"inventoryItemRepository.FindByUserIDAndID",
		func(ctx context.Context) (*dto.InventoryItemDTO, error) {
			return i.inventoryItemRepository.FindByUserIDAndID(ctx, user.ID, inventoryItemID)
		},
	)
	if err != nil {
		return err // item not found in inventory
	}

	err = tracer.TraceFn(
		ctx,
		"inventoryRepository.SetInventoryItemAsCurrent",
		func(ctx context.Context) error {
			return i.inventoryRepository.SetInventoryItemAsCurrent(ctx, user, item)
		},
	)
	if err != nil {
		return err // ???
	}

	return nil
}
