package domainservice

import (
	"abysscore/internal/adapters/controller/http/dto/request"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	"context"
)

type GameItemService interface {
	Create(
		ctx context.Context,
		request *request.CreateUpdateGameItem,
		performer *dto.UserDTO,
	) (*dto.GameItemDTO, error)

	FindAllPaged(
		ctx context.Context,
		page, size int,
		sortBy gameitementity.OrderBy,
	) (*dto.PaginatedResult[*dto.GameItemDTO], error)

	Update(
		ctx context.Context,
		id int,
		request *request.CreateUpdateGameItem,
		performer *dto.UserDTO,
	) error

	Delete(
		ctx context.Context,
		id int,
		performer *dto.UserDTO,
	) (*dto.GameItemDTO, error)
}
