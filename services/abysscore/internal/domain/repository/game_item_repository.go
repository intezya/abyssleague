package repositoryports

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	"context"
)

type GameItemRepository interface {
	Create(ctx context.Context, gameItem *dto.GameItemDTO) (*dto.GameItemDTO, error)
	FindAllPaged(
		ctx context.Context,
		page, size int,
		sortBy gameitementity.OrderBy,
	) (*dto.PaginatedResult[*dto.GameItemDTO], error)
	UpdateByID(ctx context.Context, id int, gameItem *dto.GameItemDTO) error
	DeleteByID(ctx context.Context, id int) (*dto.GameItemDTO, error)
}
