package config

import (
	"abysslib/dotenv"
	"abysslib/logger"
	"strings"
	"time"
)

type Config struct {
	GRPCPortStartFrom int

	HTTPPort int

	jwtSecret         string
	jwtIssuer         string
	jwtExpirationTime time.Duration

	Hubs []string
}

func (c Config) SecretKey() []byte {
	return []byte(c.jwtSecret)
}

func (c Config) Issuer() string {
	return c.jwtIssuer
}

func (c Config) ExpirationTime() time.Duration {
	return c.jwtExpirationTime
}

func Configure() *Config {
	logger.New(
		dotenv.GetEnv("ENV_TYPE", "dev") == "dev",
		"",
		dotenv.GetEnv("ENV_TYPE", "dev"),
	)

	logger.Log.Debugf("Debug mode: %t", true)

	logger.Log.Info("Configure success")

	return &Config{
		GRPCPortStartFrom: dotenv.GetEnvInt("GRPC_SERVER_PORT_START_FROM", 50051),

		HTTPPort: dotenv.GetEnvInt("HTTP_PORT", 8090),

		jwtSecret:         dotenv.GetEnv("JWT_SECRET", ""),
		jwtIssuer:         dotenv.GetEnv("JWT_ISSUER", "issuer"),
		jwtExpirationTime: time.Hour * 24,

		Hubs: strings.Split(dotenv.GetEnv("WEBSOCKET_HUBS", "main"), ","),
	}
}

func Setup() *Config {
	dotenv.LoadEnv()
	return Configure()
}
