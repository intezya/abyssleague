package config

import (
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/mail"
)

func initSMTPConfig() *mail.SMTPConfig {
	return &mail.SMTPConfig{
		Host:          getEnvString("SMTP_HOST", ""),
		Port:          getEnvInt("SMTP_PORT", 0),
		AccessKey:     getEnvString("SMTP_ACCESS_KEY", ""),
		SecretKey:     getEnvString("SMTP_SECRET_KEY", ""),
		DefaultSender: getEnvString("SMTP_DEFAULT_SENDER", ""),
	}
}
