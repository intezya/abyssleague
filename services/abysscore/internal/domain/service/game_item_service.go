package domainservice

import (
	"context"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/request"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/gameitementity"
)

type GameItemService interface {
	Create(
		ctx context.Context,
		request *request.CreateUpdateGameItem,
		performer *dto.UserDTO,
	) (*dto.GameItemDTO, error)

	FindByID(ctx context.Context, id int) (*dto.GameItemDTO, error)

	FindAllPaged(
		ctx context.Context,
		query *request.PaginationQuery[gameitementity.OrderBy],
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
	) error
}
