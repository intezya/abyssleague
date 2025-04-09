package config

import (
	"abysslib/dotenv"
	"abysslib/logger"
	"encoding/json"
	"log"
	"strings"
	"time"
)

const (
	DefaultGRPCPort      = 50051
	DefaultHTTPPort      = 8090
	DefaultJWTIssuer     = "issuer"
	DefaultJWTExpiration = 24 * time.Hour
	DefaultEnvType       = "dev"
	DefaultHub           = "main"

	DefaultLokiURL       = "http://localhost:3100/loki/api/v1/push"
	DefaultLokiBatchSize = 100
	DefaultLokiTimeout   = 5
)

type Config struct {
	GRPCPortStartFrom int
	HTTPPort          int

	jwtSecret         string
	jwtIssuer         string
	jwtExpirationTime time.Duration

	Hubs []string

	EnvType string
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

func (c Config) IsDevMode() bool {
	return c.EnvType == "dev"
}

func initLogger(envType string) {
	isDevMode := envType == "dev"

	lokiLabelsStr := dotenv.GetEnv("LOKI_LABELS", "")
	labels := make(map[string]string)

	if lokiLabelsStr != "" {
		if err := json.Unmarshal([]byte(lokiLabelsStr), &labels); err != nil {
			log.Printf("Error loading LOKI_LABELS: %v. Is it correct?", err)
		}
	}

	lokiConfig := logger.NewLokiConfig(
		dotenv.GetEnv("LOKI_ENDPOINT_URL", DefaultLokiURL),
		labels,
		dotenv.GetEnvInt("LOKI_BATCH_SIZE", DefaultLokiBatchSize),
		time.Duration(dotenv.GetEnvInt("LOKI_TIMEOUT_SECONDS", DefaultLokiTimeout))*time.Second,
	)

	logger.New(
		isDevMode,
		"",
		envType,
		lokiConfig,
	)

	logger.Log.Debugf("Debug mode: %t", isDevMode)
}

func Configure() *Config {
	envType := dotenv.GetEnv("ENV_TYPE", DefaultEnvType)

	initLogger(envType)

	config := &Config{
		GRPCPortStartFrom: dotenv.GetEnvInt("GRPC_SERVER_PORT_START_FROM", DefaultGRPCPort),
		HTTPPort:          dotenv.GetEnvInt("HTTP_PORT", DefaultHTTPPort),

		jwtSecret:         dotenv.GetEnv("JWT_SECRET", ""),
		jwtIssuer:         dotenv.GetEnv("JWT_ISSUER", DefaultJWTIssuer),
		jwtExpirationTime: DefaultJWTExpiration,

		Hubs: strings.Split(dotenv.GetEnv("WEBSOCKET_HUBS", DefaultHub), ","),

		EnvType: envType,
	}

	logger.Log.Info("Configuration loaded successfully")

	return config
}

func Setup() *Config {
	dotenv.LoadEnv()
	return Configure()
}
