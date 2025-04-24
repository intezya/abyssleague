package rediswrapper

import (
	"context"
	"github.com/intezya/pkglib/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Options    *redis.Options
	RetryDelay time.Duration
}

type ClientWrapper struct {
	Client            *redis.Client
	config            *Config
	closeConnectingCh chan struct{}
}

func NewClientWrapper(config *Config) *ClientWrapper {
	wrapper := &ClientWrapper{
		config:            config,
		closeConnectingCh: make(chan struct{}),
	}
	go wrapper.runConnecting(context.Background())

	return wrapper
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
				logger.Log.Info("Redis connected successfully!")

				return
			}

			logger.Log.Warnf("Attempt %d: Failed to connect to Redis: %v", attempt, err)

			attempt++
		}
	}
}

func (w *ClientWrapper) Close() {
	close(w.closeConnectingCh)

	if w.Client != nil {
		_ = w.Client.Close()
	}
}
