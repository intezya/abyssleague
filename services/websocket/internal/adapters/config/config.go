package config

import (
	"abysslib/dotenv"
	"abysslib/logger"
	"time"
)

type Config struct {
	MainGRPCHost string
	MainGRPCPort int

	DraftGRPCHost string
	DraftGRPCPort int

	HTTPPort int

	jwtSecret         string
	jwtIssuer         string
	jwtExpirationTime time.Duration
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
		false,
		"",
		dotenv.GetEnv("ENV_TYPE", "dev"),
	)

	logger.Log.Debugf("Debug mode: %t", true)

	logger.Log.Info("Configure success")

	return &Config{
		MainGRPCHost: dotenv.GetEnv("MAIN_GRPC_HOST", "localhost"),
		MainGRPCPort: dotenv.GetEnvInt("MAIN_GRPC_PORT", 50051),

		DraftGRPCHost: dotenv.GetEnv("DRAFT_GRPC_HOST", "localhost"),
		DraftGRPCPort: dotenv.GetEnvInt("DRAFT_GRPC_PORT", 50052),

		HTTPPort: dotenv.GetEnvInt("HTTP_PORT", 8090),

		jwtSecret:         dotenv.GetEnv("JWT_SECRET", ""),
		jwtIssuer:         dotenv.GetEnv("JWT_ISSUER", "issuer"),
		jwtExpirationTime: time.Hour * 24,
	}
}
