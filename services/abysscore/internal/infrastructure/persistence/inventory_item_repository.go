package persistence

import (
	"abysscore/internal/adapters/mapper"
	repositoryerrors "abysscore/internal/common/errors/repository"
	"abysscore/internal/domain/dto"
	"abysscore/internal/infrastructure/ent"
	"abysscore/internal/infrastructure/ent/inventoryitem"
	"context"
	"github.com/intezya/pkglib/itertools"
)

type InventoryItemRepository struct {
	client *ent.Client
}

func NewInventoryItemRepository(client *ent.Client) *InventoryItemRepository {
	return &InventoryItemRepository{client: client}
}

func (r *InventoryItemRepository) Create(
	ctx context.Context,
	inventoryItem *dto.CreateInventoryItemDTO,
) (*dto.InventoryItemDTO, error) {
	item, err := r.client.InventoryItem.
		Create().
		SetItemID(inventoryItem.ItemID).
		SetUserID(inventoryItem.UserID).
		SetReceivedFromID(inventoryItem.ReceivedFromID).
		Save(ctx)

	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToInventoryItemDTOFromEnt(item), nil
}

func (r *InventoryItemRepository) FindByUserID(ctx context.Context, userID int) ([]*dto.InventoryItemDTO, error) {
	result, err := r.client.InventoryItem.
		Query().
		Where(inventoryitem.UserIDEQ(userID)).
		WithItem().
		All(ctx)

	if err != nil {
		return nil, r.handleQueryError(err)
	}

	mapped := itertools.Map(result, func(v *ent.InventoryItem) *dto.InventoryItemDTO {
		return mapper.ToInventoryItemDTOFromEnt(v)
	})

	return mapped, nil
}

func (r *InventoryItemRepository) ExistsByUserIDAndID(ctx context.Context, userID, id int) bool {
	exists, err := r.client.InventoryItem.Query().
		Where(
			inventoryitem.IDEQ(id),
			inventoryitem.UserIDEQ(userID),
		).
		Exist(ctx)

	return err == nil && exists
}

func (r *InventoryItemRepository) Delete(ctx context.Context, inventoryItemID int) error {
	err := r.client.InventoryItem.
		DeleteOneID(inventoryItemID).
		Exec(ctx)

	if err != nil {
		return r.handleQueryError(err)
	}

	return nil
}

// handleQueryError transforms Ent query errors into domain-specific errors
func (r *InventoryItemRepository) handleQueryError(err error) error {
	if err == nil {
		return nil
	}

	if ent.IsNotFound(err) {
		return repositoryerrors.WrapItemNotFoundInInventory(err)
	}

	return repositoryerrors.WrapUnexpectedError(err)
}
