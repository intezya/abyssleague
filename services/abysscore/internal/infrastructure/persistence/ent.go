package persistence

import (
	"abysscore/internal/infrastructure/ent"
	"context"
	"fmt"
	"github.com/intezya/pkglib/logger"
	"time"
)

type EntConfig struct {
	driverName string
	source     string
	maxRetries int
	retryDelay time.Duration
}

func NewEntConfig(
	driverName string,
	source string,
	maxRetries int,
	retryDelay time.Duration,
) *EntConfig {
	return &EntConfig{
		driverName: driverName,
		source:     source,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

func SetupEnt(config *EntConfig) *ent.Client {
	maxRetries := gt0(config.maxRetries, 5)
	retryDelay := gt0(config.retryDelay, 2*time.Second)

	var entClient *ent.Client
	var err error

	// Retry connecting to the database if it fails
	for attempt := 1; attempt <= maxRetries; attempt++ {
		entClient, err = ent.Open(config.driverName, config.source)

		if err == nil {
			logger.Log.Infof("Database connection succeeded on attempt %d", attempt)
			break
		}

		logger.Log.Warnf("Attempt %d of %d: Failed to connect to database: %v", attempt, maxRetries, err)

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if entClient == nil {
		panic(fmt.Errorf("all attempts to connect to database failed"))
	}

	err = entClient.Schema.Create(context.Background())

	if err != nil {
		panic(fmt.Errorf("failed to create schema: %v", err))
	}

	return entClient
}

func gt0[T int | time.Duration](value T, fallback T) T {
	if value <= 0 {
		return fallback
	}
	return value
}
