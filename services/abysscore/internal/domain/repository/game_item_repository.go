package repositoryports

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	"abysscore/internal/domain/entity/types"
	"context"
)

type GameItemRepository interface {
	Create(ctx context.Context, gameItem *dto.GameItemDTO) (*dto.GameItemDTO, error)

	FindByID(ctx context.Context, id int) (*dto.GameItemDTO, error)

	FindAllPaged(
		ctx context.Context,
		page, size int,
		orderBy gameitementity.OrderBy,
		orderType types.OrderType,
	) (*dto.PaginatedResult[*dto.GameItemDTO], error)

	UpdateByID(ctx context.Context, id int, gameItem *dto.GameItemDTO) error
	DeleteByID(ctx context.Context, id int) error
}
