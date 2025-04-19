package persistence

import (
	repositoryerrors "abysscore/internal/common/errors/repository"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	"abysscore/internal/infrastructure/ent"
	"context"
	"github.com/intezya/pkglib/itertools"
)

type GameItemRepository struct {
	client *ent.Client
}

func (g *GameItemRepository) Create(ctx context.Context, gameItem *dto.GameItemDTO) (*dto.GameItemDTO, error) {
	result, err := g.client.GameItem.
		Create().
		SetName(gameItem.Name).
		SetCollection(gameItem.Collection).
		SetType(gameItem.Type).
		SetRarity(gameItem.Rarity).
		Save(ctx)

	if err != nil {
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	return dto.MapToGameItemDTOFromEnt(result), nil
}

func (g *GameItemRepository) UpdateByID(ctx context.Context, id int, gameItem *dto.GameItemDTO) error {
	_, err := g.client.GameItem.
		UpdateOneID(id).
		SetName(gameItem.Name).
		SetCollection(gameItem.Collection).
		SetType(gameItem.Type).
		SetRarity(gameItem.Rarity).
		Save(ctx)

	if err != nil {
		return repositoryerrors.WrapUnexpectedError(err)
	}

	return nil
}

func (g *GameItemRepository) DeleteByID(ctx context.Context, id int) (*dto.GameItemDTO, error) {
	return withTx(ctx, g.client, func(tx *ent.Tx) (*dto.GameItemDTO, error) {
		gameItemEnt, err := tx.GameItem.Get(ctx, id)

		if err != nil {
			if ent.IsNotFound(err) {
				return nil, repositoryerrors.WrapErrGameItemNotFound(err)
			}
			return nil, repositoryerrors.WrapUnexpectedError(err)
		}

		err = tx.GameItem.DeleteOneID(id).Exec(ctx)

		if err != nil {
			// There is cannot be not found
			return nil, repositoryerrors.WrapUnexpectedError(err)
		}

		return dto.MapToGameItemDTOFromEnt(gameItemEnt), nil
	})
}

func (g *GameItemRepository) FindAllPaged(
	ctx context.Context,
	page, size int,
	orderBy gameitementity.OrderBy,
) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
	page = getValidPage(page)
	size = getValidSize(size)
	offset := countOffset(page, size)

	total, err := g.client.GameItem.Query().Count(ctx)

	if err != nil {
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	totalPages := getTotalPages(total, size)

	items, err := g.client.GameItem.
		Query().
		Limit(size).
		Offset(offset).
		Order(orderBy.ToOrderOption()).
		All(ctx)

	if err != nil {
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	mappedItems := itertools.Map(func(item *ent.GameItem) *dto.GameItemDTO {
		return dto.MapToGameItemDTOFromEnt(item)
	}, items)

	return &dto.PaginatedResult[*dto.GameItemDTO]{
		Data:       mappedItems,
		Page:       page,
		Size:       size,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}
