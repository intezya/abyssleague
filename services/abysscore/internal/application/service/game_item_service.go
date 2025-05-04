package applicationservice

import (
	"context"

	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/http/dto/request"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/gameitementity"
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
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
	ctx, span := tracer.StartSpan(ctx, "GameItemService.Create")
	defer span.End()

	// TODO: log performer action

	result, err := g.gameItemRepository.Create(ctx, request.ToDTO())
	if err != nil {
		return nil, err // ???
	}

	return result, nil
}

func (g *GameItemService) FindByID(ctx context.Context, id int) (*dto.GameItemDTO, error) {
	ctx, span := tracer.StartSpan(ctx, "GameItemService.FindByID")
	defer span.End()

	result, err := g.gameItemRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err // not found
	}

	return result, nil
}

func (g *GameItemService) FindAllPaged(
	ctx context.Context,
	query *request.PaginationQuery[gameitementity.OrderBy],
) (*dto.PaginatedResult[*dto.GameItemDTO], error) {
	ctx, span := tracer.StartSpan(ctx, "GameItemService.FindAllPaged")
	defer span.End()

	result, err := g.gameItemRepository.FindAllPaged(
		ctx,
		query.Page,
		query.Size,
		query.OrderBy,
		query.OrderType,
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
	ctx, span := tracer.StartSpan(ctx, "GameItemService.Update")
	defer span.End()

	// TODO: log performer action

	err := g.gameItemRepository.UpdateByID(ctx, id, request.ToDTO())
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
	ctx, span := tracer.StartSpan(ctx, "GameItemService.Delete")
	defer span.End()

	// TODO: log performer action

	err := g.gameItemRepository.DeleteByID(ctx, id)
	if err != nil {
		return err // not found
	}

	return nil
}
