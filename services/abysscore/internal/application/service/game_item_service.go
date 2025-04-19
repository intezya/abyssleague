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

func NewGameItemService(repository repositoryports.GameItemRepository) *GameItemService {
	return &GameItemService{repository: repository}
}

func (g *GameItemService) Create(
	ctx context.Context,
	request *request.CreateUpdateGameItem,
	performer *dto.UserDTO,
) (*dto.GameItemDTO, error) {
	// TODO: log performer action

	return g.repository.Create(ctx, request.ToDTO())
}

func (g *GameItemService) FindByID(ctx context.Context, id int) (*dto.GameItemDTO, error) {
	return g.repository.FindByID(ctx, id)
}

func (g *GameItemService) FindAllPaged(
	ctx context.Context,
	query *request.PaginationQuery[gameitementity.OrderBy],
) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
	return g.repository.FindAllPaged(ctx, query.Page, query.Size, query.OrderBy, query.OrderType)
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
) error {
	// TODO: log performer action

	return g.repository.DeleteByID(ctx, id)
}
