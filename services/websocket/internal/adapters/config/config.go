package config

import (
	"abysslib/dotenv"
	"abysslib/itertools"
	"abysslib/logger"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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
	GRPCPorts []int
	HTTPPort  int

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

	grpcPorts := itertools.Map(
		func(s string) int {
			if i, err := strconv.Atoi(s); err != nil {
				panic(fmt.Sprintf("Error parsing GRPC_SERVER_PORTS: %s", err))
			} else {
				return i
			}
		},
		strings.Split(dotenv.GetEnv("GRPC_SERVER_PORTS", string(int32(DefaultGRPCPort))), ","),
	)
	websocketHubs := strings.Split(dotenv.GetEnv("WEBSOCKET_HUBS", DefaultHub), ",")

	if len(grpcPorts) != len(websocketHubs) {
		panic("GRPC_SERVER_PORTS and WEBSOCKET_HUBS must have the same number of elements")
	}

	config := &Config{
		GRPCPorts: grpcPorts,
		HTTPPort:  dotenv.GetEnvInt("HTTP_PORT", DefaultHTTPPort),

		jwtSecret:         dotenv.GetEnv("JWT_SECRET", ""),
		jwtIssuer:         dotenv.GetEnv("JWT_ISSUER", DefaultJWTIssuer),
		jwtExpirationTime: DefaultJWTExpiration,

		Hubs: websocketHubs,

		EnvType: envType,
	}

	logger.Log.Info("Configuration loaded successfully")

	return config
}

func Setup() *Config {
	dotenv.LoadEnv()
	return Configure()
}
