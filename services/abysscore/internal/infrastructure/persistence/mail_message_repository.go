package persistence

import (
	"context"
	"errors"
	"fmt"
	repositoryerrors "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/repository"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/mailmessage"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/pkglib/logger"
	"strings"
	"time"
)

type entry struct {
	key           string
	expireSeconds int
	data          interface{}
}

var (
	ErrClientNotReady = errors.New("client not ready")
)

type MailMessageRepository struct {
	redisClient *rediswrapper.ClientWrapper

	queue chan entry
}

func NewMailMessageRepository(redisClient *rediswrapper.ClientWrapper) *MailMessageRepository {
	const saveDataBusBuffer = 256

	repository := &MailMessageRepository{
		redisClient: redisClient,
		queue:       make(chan entry, saveDataBusBuffer),
	}

	go repository.saveWorker()

	return repository
}

func (s *MailMessageRepository) saveWorker() {
	for data := range s.queue {
		if s.redisClient.Client != nil {
			ctx := context.Background()
			err := s.redisClient.Client.Set(
				ctx,
				data.key,
				data.data,
				time.Duration(data.expireSeconds)*time.Second,
			).Err()

			if err != nil {
				logger.Log.Debug("failed to save data in cache. " + err.Error())
			}
		}
	}
}

func (s *MailMessageRepository) linkMailCodeKey(userID int) string {
	const key = "LinkMailCodeData"

	return fmt.Sprintf("%s:%d", key, userID)
}

func (s *MailMessageRepository) SaveLinkMailCodeData(ctx context.Context, message *mailmessage.LinkEmailCodeData) {
	const linkMailCodeExpireSeconds = 120

	key := s.linkMailCodeKey(message.UserID)

	data := entry{
		key:           key,
		expireSeconds: linkMailCodeExpireSeconds,
		data:          message,
	}

	s.queue <- data
}

func (s *MailMessageRepository) GetLinkMailCodeData(
	ctx context.Context,
	receiverID int,
) (*mailmessage.LinkEmailCodeData, error) {
	if s.redisClient.Client != nil {
		key := s.linkMailCodeKey(receiverID)

		result := &mailmessage.LinkEmailCodeData{}

		err := s.redisClient.Client.Get(ctx, key).Scan(result)

		if err != nil {
			return nil, s.handleNotFoundOrUnexpected(err)
		}

		return result, nil
	}

	return nil, ErrClientNotReady
}

func (s *MailMessageRepository) handleNotFoundOrUnexpected(err error) error {
	if strings.Contains(err.Error(), "redis: nil") {
		return repositoryerrors.WrapMailDataNotFound(err)
	}

	return repositoryerrors.WrapUnexpectedError(err)
}
