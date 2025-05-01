package rediswrapper

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Options    *redis.Options
	RetryDelay time.Duration
}

type ClientWrapper struct {
	Client            *redis.Client
	config            *Config
	closeConnectingCh chan struct{}
	logger            Logger
}

type Logger interface {
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Warnf(template string, args ...interface{})
}

func NewClientWrapper(config *Config, logger Logger) *ClientWrapper {
	wrapper := &ClientWrapper{
		Client:            nil,
		config:            config,
		closeConnectingCh: make(chan struct{}),
		logger:            logger,
	}

	go wrapper.runConnecting(context.Background())

	return wrapper
}

func (w *ClientWrapper) Close() {
	close(w.closeConnectingCh)

	if w.Client != nil {
		_ = w.Client.Close()
	}
}

func (w *ClientWrapper) runConnecting(ctx context.Context) {
	attempt := 1

	for {
		select {
		case <-w.closeConnectingCh:
			return
		default:
			if w.Client != nil {
				_ = w.Client.Close()
			}

			time.Sleep(w.config.RetryDelay)

			w.Client = redis.NewClient(w.config.Options)

			_, err := w.Client.Ping(ctx).Result()
			if err == nil {
				w.logger.Infoln("Redis connected successfully!")

				return
			}

			w.logger.Warnf("Attempt %d: Failed to connect to Redis: %v", attempt, err)

			attempt++
		}
	}
}
