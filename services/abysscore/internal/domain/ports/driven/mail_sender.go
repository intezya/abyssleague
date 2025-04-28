package drivenports

import (
	"context"
	"errors"
	"regexp"
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
	Send(ctx context.Context, receiver Email, data interface{}) error
}
