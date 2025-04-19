package persistence

import (
	repositoryports "abysscore/internal/domain/repository"
	"abysscore/internal/infrastructure/ent"
)

type DependencyProvider struct {
	client *ent.Client

	UserRepository     repositoryports.UserRepository
	GameItemRepository repositoryports.GameItemRepository
}

func NewDependencyProvider(client *ent.Client) *DependencyProvider {
	return &DependencyProvider{
		client: client,

		UserRepository:     NewUserRepository(client),
		GameItemRepository: NewGameItemRepository(client),
	}
}
