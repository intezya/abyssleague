package applicationservice

import (
	"abysscore/internal/adapters/controller/http/dto/request"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	repositoryports "abysscore/internal/domain/repository"
	"context"
)

type GameItemService struct {
	repository repositoryports.GameItemRepository
}

func (g *GameItemService) Create(
	ctx context.Context,
	request *request.CreateUpdateGameItem,
	performer *dto.UserDTO,
) (*dto.GameItemDTO, error) {
	// TODO: log performer action

	return g.repository.Create(ctx, request.ToDTO())
}

func (g *GameItemService) FindAllPaged(
	ctx context.Context,
	page, size int,
	sortBy gameitementity.OrderBy,
) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
	return g.repository.FindAllPaged(ctx, page, size, sortBy)
}

func (g *GameItemService) Update(
	ctx context.Context,
	id int,
	request *request.CreateUpdateGameItem,
	performer *dto.UserDTO,
) error {
	// TODO: log performer action

	return g.repository.UpdateByID(ctx, id, request.ToDTO())
}

func (g *GameItemService) Delete(
	ctx context.Context,
	id int,
	performer *dto.UserDTO,
) (*dto.GameItemDTO, error) {
	// TODO: log performer action

	return g.repository.DeleteByID(ctx, id)
}
