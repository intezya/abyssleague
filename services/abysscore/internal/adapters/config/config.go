package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/factory"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/auth"
	"github.com/intezya/pkglib/itertools"
	"github.com/intezya/pkglib/logger"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"strings"

	// ...
	"github.com/joho/godotenv"
	"time"
)

const (
	defaultServerPort             = 8080
	defaultJWTExpiration          = 24 * time.Hour
	defaultRedisRetryDelay        = 2 * time.Second
	defaultDBMaxRetries           = 5
	defaultDBRetryDelay           = 2 * time.Second
	defaultLokiBatchSize          = 100
	defaultLokiMaxWait            = 5 * time.Second
	defaultLokiTimeout            = 10 * time.Second
	defaultLokiRetryCount         = 3
	defaultLokiRetryWait          = 1 * time.Second
	defaultRateLimitLoginTime     = 10 * time.Second
	defaultRateLimitLoginRequests = 3
	defaultRateLimitDefaultTime   = 1 * time.Second
	defaultRateLimitDefaultReqs   = 4
	defaultMetricsPort            = 2112
	splitPairParts                = 2
	defaultSlowRequestThresholdMs = 300
)

type Config struct {
	FiberHealthCheckConfig healthcheck.Config
	FiberRequestIDConfig   requestid.Config
	IsDebug                bool
	EnvType                string

	RateLimitConfig  *RateLimitConfig
	LoggerConfig     *LoggerConfig
	RedisConfig      *rediswrapper.Config
	EntConfig        *persistence.EntConfig
	JWTConfiguration *auth.JWTConfiguration
	TracerConfig     *tracer.Config
	GRPCConfig       *factory.GRPCConfig

	ServerPort  int
	MetricsPort int

	SlowRequestThresholdMs int
}

type JWTConfiguration struct {
	SecretKey      []byte
	Issuer         string
	ExpirationTime time.Duration
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	envType := getEnvString("ENV_TYPE", "prod")

	config := &Config{
		FiberHealthCheckConfig: healthcheck.ConfigDefault,
		FiberRequestIDConfig:   requestid.ConfigDefault,
		IsDebug:                getEnvBool("DEBUG", false),
		EnvType:                envType,

		RateLimitConfig: &RateLimitConfig{
			LoginRateLimitKey:    getEnvString("RATE_LIMIT_LOGIN_KEY", "login_attempts_rate_limit:"),
			LoginRateLimitTime:   getEnvDuration("RATE_LIMIT_LOGIN_TIME", defaultRateLimitLoginTime),
			LoginRateLimit:       getEnvInt("RATE_LIMIT_LOGIN_MAX_REQUESTS", defaultRateLimitLoginRequests),
			DefaultRateLimitKey:  getEnvString("RATE_LIMIT_DEFAULT_KEY", "rate_limit:"),
			DefaultRateLimitTime: getEnvDuration("RATE_LIMIT_DEFAULT_TIME", defaultRateLimitDefaultTime),
			DefaultRateLimit:     getEnvInt("RATE_LIMIT_DEFAULT_MAX_REQUESTS", defaultRateLimitDefaultReqs),
		},

		LoggerConfig: &LoggerConfig{
			lokiConfig: &logger.LokiConfig{
				URL:                  getEnvString("LOKI_ENDPOINT_URL", "localhost:3100"),
				Labels:               parseLokiLabels(getEnvString("LOKI_LABELS", "")),
				BatchSize:            getEnvInt("LOKI_BATCH_SIZE", defaultLokiBatchSize),
				MaxWait:              getEnvDuration("LOKI_MAX_WAIT", defaultLokiMaxWait),
				Timeout:              getEnvDuration("LOKI_TIMEOUT", defaultLokiTimeout),
				Compression:          getEnvBool("LOKI_COMPRESSION", true),
				RetryCount:           getEnvInt("LOKI_RETRY_COUNT", defaultLokiRetryCount),
				RetryWait:            getEnvDuration("LOKI_RETRY_WAIT", defaultLokiRetryWait),
				SuppressSinkWarnings: getEnvBool("LOKI_SUPPRESS_WARNINGS", false),
			},
		},

		RedisConfig: &rediswrapper.Config{
			Options: &redis.Options{ //nolint:exhaustruct // redis config contains TOO MANY options
				Addr:     getEnvString("REDIS_ADDR", "localhost:6379"),
				Password: getEnvString("REDIS_PASSWORD", ""),
				DB:       getEnvInt("REDIS_DB", 0),
			},
			RetryDelay: getEnvDuration("REDIS_RETRY_DELAY", defaultRedisRetryDelay),
		},

		EntConfig: persistence.NewEntConfig(
			getEnvString("DB_DRIVER", "postgres"),
			buildDBConnectionString(),
			getEnvInt("DB_MAX_RETRIES", defaultDBMaxRetries),
			getEnvDuration("DB_RETRY_DELAY", defaultDBRetryDelay),
		),

		JWTConfiguration: auth.NewJWTConfiguration(
			getEnvString("JWT_SECRET", "default-secret-key"),
			getEnvString("JWT_ISSUER", "com.intezya.abyssleague.auth"),
			getEnvDuration("JWT_EXPIRATION_TIME", defaultJWTExpiration),
		),

		TracerConfig: &tracer.Config{
			Endpoint:           getEnvString("TRACER_ENDPOINT", "localhost:4317"),
			ServiceName:        getEnvString("TRACER_SERVICE_NAME", ""),
			ServiceVersion:     getEnvString("TRACER_SERVICE_VERSION", "v1"),
			ServiceEnvironment: envType,
		},

		GRPCConfig: &factory.GRPCConfig{
			WebsocketApiGatewayHost:  getEnvString("WEBSOCKET_API_GATEWAY_HOST", ""),
			WebsocketApiGatewayPorts: getAndParseGrpcPorts(),
		},

		ServerPort:             getEnvInt("SERVER_PORT", defaultServerPort),
		MetricsPort:            getEnvInt("METRICS_PORT", defaultMetricsPort),
		SlowRequestThresholdMs: defaultSlowRequestThresholdMs,
	}

	config.LoggerConfig.lokiConfig.Labels["environment"] = config.EnvType

	config.FiberRequestIDConfig.ContextKey = getEnvString("REQUEST_ID_KEY", "requestid")
	config.FiberHealthCheckConfig.LivenessEndpoint = getEnvString("LIVENESS_ENDPOINT", "/live")
	config.FiberHealthCheckConfig.ReadinessEndpoint = getEnvString("READINESS_ENDPOINT", "/health")

	return config
}

func buildDBConnectionString() string {
	if os.Getenv("DB_URL") != "" {
		return os.Getenv("DB_URL")
	}

	host := getEnvString("DB_HOST", "localhost")
	port := getEnvString("DB_PORT", "5432")
	user := getEnvString("DB_USER", "postgres")
	password := getEnvString("DB_PASSWORD", "")
	dbname := getEnvString("DB_NAME", "postgres")
	sslmode := getEnvString("DB_SSL_MODE", "disable")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

func parseLokiLabels(labelsStr string) map[string]string {
	labels := make(map[string]string)
	if labelsStr == "" {
		return labels
	}

	pairs := strings.Split(labelsStr, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", splitPairParts)
		if len(parts) == splitPairParts {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			labels[key] = value
		}
	}

	return labels
}

func getAndParseGrpcPorts() []int {
	return itertools.Map(
		strings.Split(getEnvString("WEBSOCKET_API_GATEWAY_PORTS", ""), ","),
		func(s string) int {
			serverPort, err := strconv.Atoi(s)

			if err != nil {
				panic(fmt.Sprintf("Error parsing GRPC_SERVER_PORTS: %s", err))
			}

			return serverPort
		},
	)
}
