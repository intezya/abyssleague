package persistence

import (
	repositoryports "github.com/intezya/abyssleague/services/abysscore/internal/domain/repository"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
)

type DependencyProvider struct {
	client *ent.Client

	UserRepository          repositoryports.UserRepository
	GameItemRepository      repositoryports.GameItemRepository
	InventoryItemRepository repositoryports.InventoryItemRepository
	MailMessageRepository   repositoryports.MailMessageRepository
}

func NewDependencyProvider(client *ent.Client, redisClient *rediswrapper.ClientWrapper) *DependencyProvider {
	return &DependencyProvider{
		client: client,

		UserRepository:          NewUserRepository(client),
		GameItemRepository:      NewGameItemRepository(client),
		InventoryItemRepository: NewInventoryItemRepository(client),
		MailMessageRepository:   NewMailMessageRepository(redisClient),
	}
}
