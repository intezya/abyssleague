package repositoryports

import (
	"context"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/mailmessage"
)

type MailMessageRepository interface {
	SaveLinkMailCodeData(ctx context.Context, message *mailmessage.LinkEmailCodeData, expireMinutes int)

	GetLinkMailCodeData(
		ctx context.Context,
		receiverID int,
	) (*mailmessage.LinkEmailCodeData, error)
}
