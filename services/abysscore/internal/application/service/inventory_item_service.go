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
	ctx, span := tracer.StartSpan(ctx, "InventoryItemService.GrantToUserByAdmin")
	defer span.End()

	createRequest := &dto.CreateInventoryItemDTO{
		UserID:         userID,
		ItemID:         itemID,
		ReceivedFromID: performer.ID,
	}

	result, err := i.inventoryItemRepository.Create(ctx, createRequest)
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
	ctx, span := tracer.StartSpan(ctx, "InventoryItemService.FindAllByUserID")
	defer span.End()

	items, err := i.inventoryItemRepository.FindByUserID(ctx, userID)
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
	ctx, span := tracer.StartSpan(ctx, "InventoryItemService.RevokeByAdmin")
	defer span.End()

	result, err := i.inventoryItemRepository.DeleteByUserIDAndID(ctx, userID, inventoryItemID)
	if err != nil {
		return err
	}

	i.eventPublisher.Publish(event.NewInventoryItemRevokedEvent(
		optional.EmptyOptional[string](),
		optional.New(performer),
		userID,
		result,
	))

	return nil
}

func (i *InventoryItemService) SetInventoryItemAsCurrent(
	ctx context.Context,
	user *dto.UserDTO,
	inventoryItemID int,
) error {
	ctx, span := tracer.StartSpan(ctx, "InventoryItemService.SetInventoryItemAsCurrent")
	defer span.End()

	item, err := i.inventoryItemRepository.FindByUserIDAndID(ctx, user.ID, inventoryItemID)
	if err != nil {
		return err // item not found in inventory
	}

	err = i.inventoryRepository.SetInventoryItemAsCurrent(ctx, user, item)
	if err != nil {
		return err // ???
	}

	return nil
}
