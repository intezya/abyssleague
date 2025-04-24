package persistence

import (
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
)

type DependencyProvider struct {
	client *ent.Client

	UserRepository          repositoryports.UserRepository
	GameItemRepository      repositoryports.GameItemRepository
	InventoryItemRepository repositoryports.InventoryItemRepository
}

func NewDependencyProvider(client *ent.Client) *DependencyProvider {
	return &DependencyProvider{
		client: client,

		UserRepository:          NewUserRepository(client),
		GameItemRepository:      NewGameItemRepository(client),
		InventoryItemRepository: NewInventoryItemRepository(client),
	}
}
