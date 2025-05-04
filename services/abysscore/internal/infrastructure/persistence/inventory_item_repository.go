package persistence

import (
	"context"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/mapper"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/inventoryitem"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
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
	ctx, span := tracer.StartSpan(ctx, "InventoryItemRepository.Create")
	defer span.End()

	result, err := withTxResult(ctx, r.client, func(tx *ent.Tx) (*ent.InventoryItem, error) {
		created, err := r.client.InventoryItem.
			Create().
			SetItemID(inventoryItem.ItemID).
			SetUserID(inventoryItem.UserID).
			SetReceivedFromID(inventoryItem.ReceivedFromID).
			Save(ctx)
		if err != nil {
			return nil, err
		}

		result, err := r.client.InventoryItem.
			Query().
			Where(inventoryitem.IDEQ(created.ID)).
			WithItem().
			Only(ctx)
		if err != nil {
			return nil, err
		}

		return result, nil
	})
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToInventoryItemDTOFromEnt(result), nil
}

func (r *InventoryItemRepository) FindByUserID(
	ctx context.Context,
	userID int,
) ([]*dto.InventoryItemDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "InventoryItemRepository.FindByUserID")
	defer span.End()

	inventoryItems, err := r.client.InventoryItem.
		Query().
		Where(inventoryitem.UserIDEQ(userID)).
		WithItem().
		All(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	mapped := itertools.Map(inventoryItems, mapper.ToInventoryItemDTOFromEnt)

	return mapped, nil
}

func (r *InventoryItemRepository) FindByUserIDAndID(
	ctx context.Context,
	userID, id int,
) (*dto.InventoryItemDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "InventoryItemRepository.FindByUserIDAndID")
	defer span.End()

	inventoryItem, err := r.client.InventoryItem.
		Query().
		Where(
			inventoryitem.IDEQ(id),
			inventoryitem.UserIDEQ(userID),
		).
		WithItem().
		First(ctx)
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToInventoryItemDTOFromEnt(inventoryItem), nil
}

func (r *InventoryItemRepository) DeleteByUserIDAndID(
	ctx context.Context,
	userID, id int,
) (*dto.InventoryItemDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "InventoryItemRepository.DeleteByUserIDAndID")
	defer span.End()

	result, err := withTxResult(ctx, r.client, func(tx *ent.Tx) (*ent.InventoryItem, error) {
		item, err := tx.InventoryItem.
			Query().
			Where(inventoryitem.IDEQ(id), inventoryitem.UserIDEQ(userID)).
			WithItem().
			First(ctx)
		if err != nil {
			return nil, err
		}

		err = tx.InventoryItem.DeleteOneID(item.ID).Exec(ctx)
		if err != nil {
			return nil, err
		}

		return item, nil
	})
	if err != nil {
		return nil, r.handleQueryError(err)
	}

	return mapper.ToInventoryItemDTOFromEnt(result), nil
}

// handleQueryError transforms Ent query errors into domain-specific errors.
func (r *InventoryItemRepository) handleQueryError(err error) error {
	if err == nil {
		return nil
	}

	if ent.IsNotFound(err) {
		return apperrors.WrapItemNotFoundInInventory(err)
	}

	return apperrors.WrapUnexpectedError(err)
}
