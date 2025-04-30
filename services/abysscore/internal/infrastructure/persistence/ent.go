package persistence

import (
	"context"
	"errors"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent/migrate"
	"time"

	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/ent"
	"github.com/intezya/pkglib/logger"
)

const (
	defaultEntReconnectMaxRetries = 5
	defaultEntReconnectDelay      = 2 * time.Second
)

var errAllConnectionAttemptsFailed = errors.New("all attempts to connect to database failed")

type EntConfig struct {
	driverName string
	source     string
	maxRetries int
	retryDelay time.Duration
	debug      bool
}

func NewEntConfig(
	driverName string,
	source string,
	maxRetries int,
	retryDelay time.Duration,
	debug bool,
) *EntConfig {
	return &EntConfig{
		driverName: driverName,
		source:     source,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		debug:      debug,
	}
}

func SetupEnt(config *EntConfig) *ent.Client {
	maxRetries := gt0(config.maxRetries, defaultEntReconnectMaxRetries)
	retryDelay := gt0(config.retryDelay, defaultEntReconnectDelay)

	entClient, err := ent.Open(config.driverName, config.source)

	if err != nil {
		logger.Log.Fatal(err) // invalid driver
	}

	if config.debug {
		entClient = entClient.Debug()
	}

	// Retry connecting to the database if it fails
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = entClient.Schema.Create(
			context.Background(),
			migrate.WithDropIndex(true),
			migrate.WithDropColumn(true),
		)
		if err == nil {
			logger.Log.Infof("Database migrations runned success on attempt %d", attempt)

			break
		}

		logger.Log.Warnf(
			"Attempt %d of %d: Failed to run migrations for database: %v",
			attempt,
			maxRetries,
			err,
		)

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		logger.Log.Fatalf("failed to create schema (all attempts are over)")
	}

	return entClient
}

func gt0[T int | time.Duration](value T, fallback T) T {
	if value <= 0 {
		return fallback
	}

	return value
}
