package clients

import (
	"fmt"
	"time"
)

// Server indices for websocket servers.
const (
	MainWebsocketServerIdx  = 0
	DraftWebsocketServerIdx = 1
)

const (
	DefaultRequestTimeout         = 2 * time.Second
	DefaultWebsocketMessagingHost = "websocket-messaging"
)

var DefaultWebsocketMessagingPorts = []int{50051, 50052}

// Config holds factory-wide configuration settings.
type Config struct {
	// Environment settings
	DevMode bool

	RequestTimeout time.Duration

	WebsocketMessagingServiceHost  string
	WebsocketMessagingServicePorts []int
}

func DefaultConfig() *Config {
	return &Config{
		DevMode:                        false,
		RequestTimeout:                 DefaultRequestTimeout,
		WebsocketMessagingServiceHost:  DefaultWebsocketMessagingHost,
		WebsocketMessagingServicePorts: DefaultWebsocketMessagingPorts,
	}
}

// MainWebsocketServerAddress returns the address of the main websocket server.
func (c *Config) MainWebsocketServerAddress() string {
	if c.WebsocketMessagingServiceHost == "" ||
		len(c.WebsocketMessagingServicePorts) <= MainWebsocketServerIdx {
		return ""
	}

	return fmt.Sprintf(
		"%s:%d",
		c.WebsocketMessagingServiceHost,
		c.WebsocketMessagingServicePorts[MainWebsocketServerIdx],
	)
}

// DraftWebsocketServerAddress returns the address of the draft websocket server.
func (c *Config) DraftWebsocketServerAddress() string {
	if c.WebsocketMessagingServiceHost == "" ||
		len(c.WebsocketMessagingServicePorts) <= DraftWebsocketServerIdx {
		return ""
	}

	return fmt.Sprintf(
		"%s:%d",
		c.WebsocketMessagingServiceHost,
		c.WebsocketMessagingServicePorts[DraftWebsocketServerIdx],
	)
}
