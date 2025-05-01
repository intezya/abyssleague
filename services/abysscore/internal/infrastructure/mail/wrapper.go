package mail

import (
	"context"
	"errors"
	"fmt"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/mailmessage"
	"gopkg.in/mail.v2"
	"net"
	"time"
)

type Logger interface {
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
}

type SMTPSender struct {
	config *SMTPConfig
	dialer *mail.Dialer
	logger Logger
}

func NewSMTPSender(config *SMTPConfig, logger Logger) *SMTPSender {
	const dialTimeout = 5 * time.Second
	const mailDevHost = "maildev"

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), dialTimeout)

	if err != nil {
		logger.Errorln("Failed to connect to SMTP server:", err)
	} else {
		logger.Infoln("Successfully established TCP connection to SMTP server")
		_ = conn.Close()
	}

	dialer := mail.NewDialer(config.Host, config.Port, config.AccessKey, config.SecretKey)
	dialer.Timeout = dialTimeout

	if config.Host == mailDevHost {
		dialer.SSL = false
		dialer.StartTLSPolicy = mail.NoStartTLS
		dialer.Auth = newCustomAuth(config.AccessKey, config.SecretKey)
	}

	return &SMTPSender{
		config: config,
		dialer: dialer,
		logger: logger,
	}
}

func (s *SMTPSender) Send(ctx context.Context, message *mailmessage.Message, receivers ...string) error {
	return s.SendS(ctx, s.config.DefaultSender, message, receivers...)
}

func (s *SMTPSender) SendS(ctx context.Context, sender string, message *mailmessage.Message, receivers ...string) error {
	s.logger.Debugln("Sending email from:", sender, "to:", receivers)

	if len(receivers) == 0 {
		return errors.New("at least one receiver email is required")
	}

	msg := mail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", receivers...)
	msg.SetHeader("Subject", message.Subject)

	if message.Mime != "" {
		msg.SetBody(message.Mime, message.Body)
	} else {
		msg.SetBody("text/html", message.Body)
	}

	done := make(chan error, 1)

	go func() {
		s.logger.Debugln("Dialing SMTP server...")
		err := s.dialer.DialAndSend(msg)
		if err != nil {
			s.logger.Warnln("Failed to send email:", err)
		} else {
			s.logger.Debugln("Email sent successfully")
		}
		done <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("mail sending canceled or timed out: %w", ctx.Err())
	case err := <-done:
		return err
	}
}
