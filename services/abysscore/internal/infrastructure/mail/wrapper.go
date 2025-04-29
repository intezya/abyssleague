package mail

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/config"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/mailmessage"
	"github.com/intezya/pkglib/logger"
	"net/smtp"
)

type SMTPClientWrapper struct {
	auth          smtp.Auth
	url           string
	defaultSender string
}

func (c *SMTPClientWrapper) Send(message *mailmessage.Message, receiver ...string) error {
	err := smtp.SendMail(c.url, c.auth, c.defaultSender, receiver, message.AsBytes())

	if err != nil {
		logger.Log.Debug("failed to send mail: ", err)
	}

	return nil
}

func (c *SMTPClientWrapper) SendS(sender string, message *mailmessage.Message, receiver ...string) error {
	err := smtp.SendMail(c.url, c.auth, sender, receiver, message.AsBytes())

	if err != nil {
		logger.Log.Debug("failed to send mail: ", err)
	}

	return nil
}

func NewSMTPClientWrapper(cfg *config.SMTPConfig) *SMTPClientWrapper {
	auth := smtp.PlainAuth("", cfg.AccessKey, cfg.SecretKey, cfg.Host)

	return &SMTPClientWrapper{
		auth: auth,
		url:  cfg.Server(),
	}
}
