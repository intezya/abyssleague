package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/controller/grpc/clients"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/mail"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/metrics/tracer"
	"github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/persistence"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/auth"
	"github.com/intezya/pkglib/itertools"
	"github.com/intezya/pkglib/logger"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type EnvType string

const (
	EnvTypeDev  EnvType = "dev"
	EnvTypeProd EnvType = "prod"
)

var (
	errLoginRateLimitKeyEmpty   = errors.New("login rate limit key cannot be empty")
	errDefaultRateLimitKeyEmpty = errors.New("default rate limit key cannot be empty")
	errLoadEnvFile              = errors.New("error loading .env file")
)

// Default configuration values.
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

// ConfigValidator represents configuration validation interface.
type ConfigValidator interface {
	Validate() error
}

// Config represents application configuration.
type Config struct {
	// Server configuration
	ServerPort             int
	SlowRequestThresholdMs int
	MetricsPort            int
	FiberHealthCheckConfig healthcheck.Config
	FiberRequestIDConfig   requestid.Config

	// Environment configuration
	IsDebug bool
	EnvType EnvType

	// Component configurations
	RateLimitConfig  *RateLimitConfig
	LoggerConfig     *LoggerConfig
	RedisConfig      *rediswrapper.Config
	EntConfig        *persistence.EntConfig
	JWTConfiguration *auth.JWTConfiguration
	TracerConfig     *tracer.Config
	GRPCConfig       *clients.Config
	SMTPConfig       *mail.SMTPConfig
}

// Validate validates the rate limit configuration.
func (c *RateLimitConfig) Validate() error {
	if c.LoginRateLimitKey == "" {
		return errLoginRateLimitKeyEmpty
	}

	if c.DefaultRateLimitKey == "" {
		return errDefaultRateLimitKeyEmpty
	}

	return nil
}

// GetLokiConfig returns the Loki logger configuration.
func (l *LoggerConfig) GetLokiConfig() *logger.LokiConfig {
	return l.lokiConfig
}

// LoadConfig loads the application configuration from environment variables
// and returns a complete Config structure.
func LoadConfig() *Config {
	// Load environment variables from .env file if exists
	if err := godotenv.Load(); err != nil {
		// Only log error, don't fail as .env file is optional
		//nolint:forbidigo // logger not initialized yet
		fmt.Printf("%v: %v\n", errLoadEnvFile, err)
	}

	envType := getEnvString("ENV_TYPE", string(EnvTypeProd))

	config := &Config{
		// Server configuration
		ServerPort:  getEnvInt("SERVER_PORT", defaultServerPort),
		MetricsPort: getEnvInt("METRICS_PORT", defaultMetricsPort),
		SlowRequestThresholdMs: getEnvInt(
			"SLOW_REQUEST_THRESHOLD_MS",
			defaultSlowRequestThresholdMs,
		),
		FiberHealthCheckConfig: healthcheck.ConfigDefault,
		FiberRequestIDConfig:   requestid.ConfigDefault,

		// Environment configuration
		IsDebug: getEnvBool("DEBUG", false),
		EnvType: EnvType(envType),

		// Component configurations
		RateLimitConfig: initRateLimitConfig(),
		LoggerConfig:    initLoggerConfig(envType),
		RedisConfig:     initRedisConfig(),
		EntConfig:       initEntConfig(),
		JWTConfiguration: auth.NewJWTConfiguration(
			getEnvString("JWT_SECRET", "default-secret-key"),
			getEnvString("JWT_ISSUER", "com.intezya.abyssleague.auth"),
			getEnvDuration("JWT_EXPIRATION_TIME", defaultJWTExpiration),
		),
		TracerConfig: initTracerConfig(envType),
		GRPCConfig:   initGRPCConfig(envType == string(EnvTypeDev)),
		SMTPConfig:   initSMTPConfig(),
	}

	// Set specific Fiber middleware configurations
	config.FiberRequestIDConfig.ContextKey = getEnvString("REQUEST_ID_KEY", "requestid")
	config.FiberHealthCheckConfig.LivenessEndpoint = getEnvString("LIVENESS_ENDPOINT", "/live")
	config.FiberHealthCheckConfig.ReadinessEndpoint = getEnvString("READINESS_ENDPOINT", "/health")

	// Validate critical configuration
	if err := validateConfig(config); err != nil {
		panic(fmt.Errorf("config validation failed: %w", err))
	}

	return config
}

// validateConfig performs validation checks on the configuration.
func validateConfig(config *Config) error {
	// Validate configurations that implement ConfigValidator
	validators := []struct {
		name      string
		validator ConfigValidator
	}{
		{"RateLimitConfig", config.RateLimitConfig},
	}

	for _, v := range validators {
		if err := v.validator.Validate(); err != nil {
			return fmt.Errorf("%s validation failed: %w", v.name, err)
		}
	}

	return nil
}

// initRateLimitConfig initializes rate limit configuration.
func initRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		// Login rate limiting
		LoginRateLimitKey:  getEnvString("RATE_LIMIT_LOGIN_KEY", "login_attempts_rate_limit:"),
		LoginRateLimitTime: getEnvDuration("RATE_LIMIT_LOGIN_TIME", defaultRateLimitLoginTime),
		LoginRateLimit: getEnvInt(
			"RATE_LIMIT_LOGIN_MAX_REQUESTS",
			defaultRateLimitLoginRequests,
		),

		// Default rate limiting
		DefaultRateLimitKey: getEnvString("RATE_LIMIT_DEFAULT_KEY", "rate_limit:"),
		DefaultRateLimitTime: getEnvDuration(
			"RATE_LIMIT_DEFAULT_TIME",
			defaultRateLimitDefaultTime,
		),
		DefaultRateLimit: getEnvInt(
			"RATE_LIMIT_DEFAULT_MAX_REQUESTS",
			defaultRateLimitDefaultReqs,
		),
	}
}

// initLoggerConfig initializes logger configuration.
func initLoggerConfig(envType string) *LoggerConfig {
	lokiConfig := &logger.LokiConfig{
		URL:                  getEnvString("LOKI_ENDPOINT_URL", "localhost:3100"),
		Labels:               parseLokiLabels(getEnvString("LOKI_LABELS", "")),
		BatchSize:            getEnvInt("LOKI_BATCH_SIZE", defaultLokiBatchSize),
		MaxWait:              getEnvDuration("LOKI_MAX_WAIT", defaultLokiMaxWait),
		Timeout:              getEnvDuration("LOKI_TIMEOUT", defaultLokiTimeout),
		Compression:          getEnvBool("LOKI_COMPRESSION", true),
		RetryCount:           getEnvInt("LOKI_RETRY_COUNT", defaultLokiRetryCount),
		RetryWait:            getEnvDuration("LOKI_RETRY_WAIT", defaultLokiRetryWait),
		SuppressSinkWarnings: getEnvBool("LOKI_SUPPRESS_WARNINGS", false),
	}

	// Add environment to labels
	if lokiConfig.Labels == nil {
		lokiConfig.Labels = make(map[string]string)
	}

	lokiConfig.Labels["environment"] = envType

	return &LoggerConfig{
		lokiConfig: lokiConfig,
	}
}

// initRedisConfig initializes Redis configuration.
func initRedisConfig() *rediswrapper.Config {
	return &rediswrapper.Config{
		Options: &redis.Options{
			Addr:     getEnvString("REDIS_ADDR", "localhost:6379"),
			Password: getEnvString("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		RetryDelay: getEnvDuration("REDIS_RETRY_DELAY", defaultRedisRetryDelay),
	}
}

// initEntConfig initializes database configuration.
func initEntConfig() *persistence.EntConfig {
	return persistence.NewEntConfig(
		getEnvString("DB_DRIVER", "postgres"),
		buildDBConnectionString(),
		getEnvInt("DB_MAX_RETRIES", defaultDBMaxRetries),
		getEnvDuration("DB_RETRY_DELAY", defaultDBRetryDelay),
		getEnvBool("DB_USE_DEBUG", false),
	)
}

// initTracerConfig initializes tracer configuration.
func initTracerConfig(envType string) *tracer.Config {
	return &tracer.Config{
		Endpoint:           getEnvString("TRACER_ENDPOINT", "localhost:4317"),
		ServiceName:        getEnvString("TRACER_SERVICE_NAME", ""),
		ServiceVersion:     getEnvString("TRACER_SERVICE_VERSION", "v1"),
		ServiceEnvironment: envType,
	}
}

// initGRPCConfig initializes GRPC configuration.
func initGRPCConfig(devMode bool) *clients.Config {
	defaultConfig := &clients.Config{
		DevMode: devMode,
		RequestTimeout: time.Millisecond * time.Duration(getEnvInt(
			"ABYSSCORE_GRPC_CLIENT_CALL_TIMEOUT_MS",
			int(clients.DefaultRequestTimeout.Milliseconds()),
		)),
		WebsocketMessagingServiceHost: getEnvString("WEBSOCKET_API_GATEWAY_HOST", ""),
		WebsocketMessagingServicePorts: parseGrpcPorts(
			getEnvString("WEBSOCKET_API_GATEWAY_PORTS", ""),
		),
	}

	return defaultConfig
}

// buildDBConnectionString creates a database connection string.
func buildDBConnectionString() string {
	// Use DB_URL if provided
	if dbURL := os.Getenv("DB_URL"); dbURL != "" {
		return dbURL
	}

	// Otherwise build from components
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

// parseLokiLabels converts a comma-separated key=value string into a map.
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

// parseGrpcPorts converts a comma-separated list of ports into a slice of integers.
func parseGrpcPorts(portsStr string) []int {
	if portsStr == "" {
		return []int{}
	}

	return itertools.Map(
		strings.Split(portsStr, ","),
		func(str string) int {
			str = strings.TrimSpace(str)
			if str == "" {
				return 0
			}

			port, err := strconv.Atoi(str)
			if err != nil {
				panic(fmt.Sprintf("Error parsing GRPC port '%s': %v", str, err))
			}

			return port
		},
	)
}
