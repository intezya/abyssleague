package config

import (
	"fmt"
	"github.com/intezya/abyssleague/services/websocket/internal/pkg/auth"
	"github.com/intezya/pkglib/configloader"
	"github.com/intezya/pkglib/itertools"
	"github.com/intezya/pkglib/logger"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

/*

| Variable               | Description                                | Default                                  |
|------------------------|--------------------------------------------|------------------------------------------|
| `HTTP_PORT`            | HTTP server port                           | 8090                                     |
| `JWT_SECRET`           | Secret key for JWT authentication          | Required                                 |
| `JWT_ISSUER`           | Issuer for JWT tokens                      | "issuer"                                 |
| `ENV_TYPE`             | Environment type (dev, prod)               | "dev"                                    |
| `WEBSOCKET_HUBS`       | Comma-separated list of available hubs     | "main"                                   |
| `GRPC_SERVER_PORTS`    | Comma-separated list of ports for each hub | 50051                                    |
| `LOKI_ENDPOINT_URL`    | URL for Loki logging                       | "http://localhost:3100/loki/api/v1/push" |
| `LOKI_LABELS`          | JSON-encoded map of labels for Loki        | {}                                       |

*/

const (
	DefaultGRPCPort      = 50051
	DefaultHTTPPort      = 8090
	DefaultJWTIssuer     = "issuer"
	DefaultJWTExpiration = 24 * time.Hour
	DefaultEnvType       = "dev"
	DefaultHub           = "main"

	DefaultLokiURL = "http://localhost:3100/loki/api/v1/push"
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

func (c Config) JwtConfiguration() *auth.JWTConfiguration {
	return auth.NewJWTConfiguration(c.jwtSecret, c.jwtIssuer, c.jwtExpirationTime)
}

func (c Config) IsDevMode() bool {
	return c.EnvType == "dev"
}

func initLogger(isDebug bool, envType string) {
	lokiLabels := parseLokiLabels(configloader.GetEnvOrFallback("LOKI_LABELS", ""))
	lokiLabels["environment"] = envType

	lokiConfig := logger.NewLokiConfig(
		configloader.GetEnvOrFallback("LOKI_ENDPOINT_URL", DefaultLokiURL),
		lokiLabels,
	)

	_, err := logger.New(
		logger.WithDebug(isDebug),
		logger.WithCaller(true),
		logger.WithEnvironment(envType),
		logger.WithLoki(lokiConfig),
	)

	if err != nil {
		log.Printf("Error initializing logger: %v", err)
	}

	logger.Log.Debugf("Debug mode: %t", isDebug)
}

func Configure() *Config {
	envType := configloader.GetEnvOrFallback("ENV_TYPE", DefaultEnvType)

	initLogger(getEnvBool("DEBUG", false), envType)

	grpcPorts := itertools.Map(
		func(s string) int {
			if i, err := strconv.Atoi(s); err != nil {
				panic(fmt.Sprintf("Error parsing GRPC_SERVER_PORTS: %s", err))
			} else {
				return i
			}
		},
		strings.Split(configloader.GetEnvOrFallback("GRPC_SERVER_PORTS", string(int32(DefaultGRPCPort))), ","),
	)

	websocketHubs := strings.Split(configloader.GetEnvOrFallback("WEBSOCKET_HUBS", DefaultHub), ",")

	if len(grpcPorts) != len(websocketHubs) {
		panic("GRPC_SERVER_PORTS and WEBSOCKET_HUBS must have the same number of elements")
	}

	jwtSecret := configloader.GetEnvOrPanic("JWT_SECRET")
	jwtIssuer := configloader.GetEnvOrFallback("JWT_ISSUER", DefaultJWTIssuer)
	jwtExpirationTime := DefaultJWTExpiration

	config := &Config{
		GRPCPorts: grpcPorts,
		HTTPPort:  configloader.GetEnvIntOrFallback("HTTP_PORT", DefaultHTTPPort),

		jwtSecret:         jwtSecret,
		jwtIssuer:         jwtIssuer,
		jwtExpirationTime: jwtExpirationTime,

		Hubs: websocketHubs,

		EnvType: envType,
	}

	logger.Log.Info("Configuration loaded successfully")

	return config
}

func Setup() *Config {
	configloader.LoadEnv()
	return Configure()
}

func parseLokiLabels(labelsStr string) map[string]string {
	labels := make(map[string]string)
	if labelsStr == "" {
		return labels
	}

	pairs := strings.Split(labelsStr, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			labels[key] = value
		}
	}
	return labels
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}
