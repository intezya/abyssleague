package drivenports

import (
	"context"
	"errors"
	"regexp"

	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/mailmessage"
)

type Email string

var (
	ErrInvalidEmail = errors.New("invalid email format")
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func NewEmail(email string) (Email, error) {
	if !emailRegex.MatchString(email) {
		return "", ErrInvalidEmail
	}

	return Email(email), nil
}

func (e Email) String() string {
	return string(e)
}

type MailSender interface {
	Send(ctx context.Context, message *mailmessage.Message, receiver ...string) error
	SendS(
		ctx context.Context,
		sender string,
		message *mailmessage.Message,
		receiver ...string,
	) error
}
