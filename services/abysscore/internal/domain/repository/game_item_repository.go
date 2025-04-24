package repositoryports

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/gameitementity"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/types"
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
