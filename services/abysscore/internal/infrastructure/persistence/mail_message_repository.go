package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"strings"
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/mailmessage"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
	"github.com/intezya/pkglib/logger"
)

type entry struct {
	key           string
	expireMinutes int
	data          interface{}
}

var ErrClientNotReady = errors.New("client not ready")

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

func (s *MailMessageRepository) SaveLinkMailCodeData(
	ctx context.Context,
	message *mailmessage.LinkEmailCodeData,
	expireMinutes int,
) {
	ctx, span := tracer.StartSpan(ctx, "MailMessageRepository.SaveLinkMailCodeData")
	defer span.End()

	key := s.linkMailCodeKey(message.UserID)

	data := entry{
		key:           key,
		expireMinutes: expireMinutes,
		data:          message,
	}

	s.queue <- data
}

func (s *MailMessageRepository) GetLinkMailCodeData(
	ctx context.Context,
	receiverID int,
) (*mailmessage.LinkEmailCodeData, error) {
	ctx, span := tracer.StartSpan(ctx, "MailMessageRepository.GetLinkMailCodeData")
	defer span.End()

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
		return apperrors.WrapMailDataNotFound(err)
	}

	return apperrors.WrapUnexpectedError(err)
}

func (s *MailMessageRepository) saveWorker() {
	for data := range s.queue {
		if s.redisClient.Client != nil {
			ctx := context.Background()

			err := s.redisClient.Client.Set(
				ctx,
				data.key,
				data.data,
				time.Duration(data.expireMinutes)*time.Minute,
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
