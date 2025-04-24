package applicationservice

import (
	"abysscore/internal/adapters/controller/http/dto/request"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/entity/gameitementity"
	repositoryports "abysscore/internal/domain/repository"
	"abysscore/internal/infrastructure/metrics/tracer"
	"context"
)

type GameItemService struct {
	gameItemRepository repositoryports.GameItemRepository
}

func NewGameItemService(repository repositoryports.GameItemRepository) *GameItemService {
	return &GameItemService{gameItemRepository: repository}
}

func (g *GameItemService) Create(
	ctx context.Context,
	request *request.CreateUpdateGameItem,
	performer *dto.UserDTO,
) (*dto.GameItemDTO, error) {
	// TODO: log performer action
	result, err := tracer.TraceFnWithResult(
		ctx,
		"gameItemRepository.Create",
		func(ctx context.Context) (*dto.GameItemDTO, error) {
			return g.gameItemRepository.Create(ctx, request.ToDTO())
		},
	)

	if err != nil {
		return nil, err // ???
	}

	return result, nil
}

func (g *GameItemService) FindByID(ctx context.Context, id int) (*dto.GameItemDTO, error) {
	result, err := tracer.TraceFnWithResult(
		ctx,
		"gameItemRepository.FindByID",
		func(ctx context.Context) (*dto.GameItemDTO, error) {
			return g.gameItemRepository.FindByID(ctx, id)
		},
	)

	if err != nil {
		return nil, err // not found
	}

	return result, nil
}

func (g *GameItemService) FindAllPaged(
	ctx context.Context,
	query *request.PaginationQuery[gameitementity.OrderBy],
) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
	result, err := tracer.TraceFnWithResult(
		ctx,
		"gameItemRepository.FindAllPaged",
		func(ctx context.Context) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
			return g.gameItemRepository.FindAllPaged(ctx, query.Page, query.Size, query.OrderBy, query.OrderType)
		},
	)

	if err != nil {
		return nil, err // ???
	}

	return result, nil
}

func (g *GameItemService) Update(
	ctx context.Context,
	id int,
	request *request.CreateUpdateGameItem,
	performer *dto.UserDTO,
) error {
	// TODO: log performer action
	err := tracer.TraceFn(ctx, "gameItemRepository.UpdateByID", func(ctx context.Context) error {
		return g.gameItemRepository.UpdateByID(ctx, id, request.ToDTO())
	})

	if err != nil {
		return err // not found
	}

	return nil
}

func (g *GameItemService) Delete(
	ctx context.Context,
	id int,
	performer *dto.UserDTO,
) error {
	// TODO: log performer action
	err := tracer.TraceFn(ctx, "gameItemRepository.DeleteByID", func(ctx context.Context) error {
		return g.gameItemRepository.DeleteByID(ctx, id)
	})

	if err != nil {
		return err // not found
	}

	return nil
}
