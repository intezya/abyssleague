package persistence

import (
	"abysscore/internal/adapters/mapper"
	repositoryerrors "abysscore/internal/common/errors/repository"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	"abysscore/internal/domain/entity/types"
	"abysscore/internal/infrastructure/ent"
	"context"
	"github.com/intezya/pkglib/itertools"
)

type GameItemRepository struct {
	client *ent.Client
}

func NewGameItemRepository(client *ent.Client) *GameItemRepository {
	return &GameItemRepository{client: client}
}

func (r *GameItemRepository) Create(ctx context.Context, gameItem *dto.GameItemDTO) (*dto.GameItemDTO, error) {
	result, err := r.client.GameItem.
		Create().
		SetName(gameItem.Name).
		SetCollection(gameItem.Collection).
		SetType(gameItem.Type).
		SetRarity(gameItem.Rarity).
		Save(ctx)

	if err != nil {
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	return mapper.ToGameItemDTOFromEnt(result), nil
}

func (r *GameItemRepository) UpdateByID(ctx context.Context, id int, gameItem *dto.GameItemDTO) error {
	_, err := r.client.GameItem.
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

func (r *GameItemRepository) DeleteByID(ctx context.Context, id int) error {
	err := r.client.GameItem.
		DeleteOneID(id).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return repositoryerrors.WrapGameItemNotFound(err)
		}

		return repositoryerrors.WrapUnexpectedError(err)
	}

	return nil
}

func (r *GameItemRepository) FindAllPaged(
	ctx context.Context,
	page, size int,
	orderBy gameitementity.OrderBy,
	orderType types.OrderType,
) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
	page = getValidPage(page)
	size = getValidSize(size)
	offset := countOffset(page, size)

	total, err := r.client.GameItem.
		Query().
		Count(ctx)

	if err != nil {
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	totalPages := getTotalPages(total, size)

	gameItems, err := r.client.GameItem.
		Query().
		Limit(size).
		Offset(offset).
		Order(orderBy.ToOrderOption(orderType)).
		All(ctx)

	if err != nil {
		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	mappedItems := itertools.Map(gameItems, func(item *ent.GameItem) *dto.GameItemDTO {
		return mapper.ToGameItemDTOFromEnt(item)
	})

	return &dto.PaginatedResult[*dto.GameItemDTO]{
		Data:       mappedItems,
		Page:       page,
		Size:       size,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}

func (r *GameItemRepository) FindByID(ctx context.Context, id int) (*dto.GameItemDTO, error) {
	result, err := r.client.GameItem.
		Get(ctx, id)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, repositoryerrors.WrapGameItemNotFound(err)
		}

		return nil, repositoryerrors.WrapUnexpectedError(err)
	}

	return mapper.ToGameItemDTOFromEnt(result), nil
}
