package domainservice

import (
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	"context"
)

type GameItemService interface {
	Create(ctx context.Context, gameItem *dto.GameItemDTO, performer *dto.UserDTO) (*dto.GameItemDTO, error)
	FindAllPaged(ctx context.Context, page, size int, sortBy gameitementity.OrderBy) (*dto.PaginatedResult, error)
	Update(ctx context.Context, gameItem *dto.GameItemDTO, performer *dto.UserDTO) (*dto.GameItemDTO, error)
	Delete(ctx context.Context, gameItem *dto.GameItemDTO, performer *dto.UserDTO) error
}
