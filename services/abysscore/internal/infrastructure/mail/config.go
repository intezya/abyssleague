package mail

import "fmt"

type SMTPConfig struct {
	Host          string
	Port          int
	AccessKey     string
	SecretKey     string
	DefaultSender string
}

func (c *SMTPConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
