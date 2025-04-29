package config

import "fmt"

type SMTPConfig struct {
	Host      string
	Port      int
	AccessKey string
	SecretKey string
}

func (c *SMTPConfig) Server() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func initSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:      getEnvString("SMTP_HOST", ""),
		Port:      getEnvInt("SMTP_PORT", 0),
		AccessKey: getEnvString("SMTP_ACCESS_KEY", ""),
		SecretKey: getEnvString("SMTP_SECRET_KEY", ""),
	}
}
