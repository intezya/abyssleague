package config

import (
	rediswrapper "abysscore/internal/infrastructure/cache/redis"
	"abysscore/internal/infrastructure/metrics/tracer"
	"abysscore/internal/infrastructure/persistence"
	"abysscore/pkg/auth"
	"errors"
	"fmt"
	"github.com/intezya/pkglib/logger"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Config struct {
	UnprotectedAuthRequests []*regexp.Regexp
	FiberHealthCheckConfig  healthcheck.Config
	FiberRequestIDConfig    requestid.Config
	IsDebug                 bool
	EnvType                 string

	RateLimitConfig  *RateLimitConfig
	LoggerConfig     *LoggerConfig
	RedisConfig      *rediswrapper.Config
	EntConfig        *persistence.EntConfig
	JWTConfiguration *auth.JWTConfiguration
	TracerConfig     *tracer.Config

	MetricsPort int

	Paths *PathConfig

	SlowRequestThresholdMs int
}

type JWTConfiguration struct {
	SecretKey      []byte
	Issuer         string
	ExpirationTime time.Duration
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	config := &Config{
		IsDebug: getEnvBool("DEBUG", false),
		EnvType: getEnvString("ENV_TYPE", "prod"),
		Paths:   &PathConfig{},
	}

	config.JWTConfiguration = auth.NewJWTConfiguration(
		getEnvString("JWT_SECRET", "default-secret-key"),
		getEnvString("JWT_ISSUER", "com.intezya.abyssleague.auth"),
		getEnvDuration("JWT_EXPIRATION_TIME", 24*time.Hour),
	)

	config.RedisConfig = &rediswrapper.Config{
		Options: &redis.Options{
			Addr:     getEnvString("REDIS_ADDR", "localhost:6379"),
			Password: getEnvString("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		RetryDelay: getEnvDuration("REDIS_RETRY_DELAY", 2*time.Second),
	}

	config.EntConfig = persistence.NewEntConfig(
		getEnvString("DB_DRIVER", "postgres"),
		buildDBConnectionString(),
		getEnvInt("DB_MAX_RETRIES", 5),
		getEnvDuration("DB_RETRY_DELAY", 2*time.Second),
	)

	// Setup Logger Configuration
	config.LoggerConfig = &LoggerConfig{
		lokiConfig: &logger.LokiConfig{
			URL:         getEnvString("LOKI_ENDPOINT_URL", "localhost:3100"),
			Labels:      parseLokiLabels(getEnvString("LOKI_LABELS", "")),
			BatchSize:   getEnvInt("LOKI_BATCH_SIZE", 100),
			MaxWait:     getEnvDuration("LOKI_MAX_WAIT", 5*time.Second),
			Timeout:     getEnvDuration("LOKI_TIMEOUT", 10*time.Second),
			Compression: getEnvBool("LOKI_COMPRESSION", true),
			RetryCount:  getEnvInt("LOKI_RETRY_COUNT", 3),
			RetryWait:   getEnvDuration("LOKI_RETRY_WAIT", 1*time.Second),
		},
	}

	config.LoggerConfig.lokiConfig.Labels["environment"] = config.EnvType

	unprotected := strings.Split(getEnvString("UNPROTECTED_AUTH_ROUTES", "/api/v1/auth/.*"), ",")

	config.UnprotectedAuthRequests = make([]*regexp.Regexp, 0, len(unprotected))
	for _, route := range unprotected {
		if route == "" {
			continue
		}
		re, err := regexp.Compile(strings.TrimSpace(route))
		if err != nil {
			panic(errors.New(fmt.Sprintf("invalid unprotected route pattern %s: %v", route, err)))
		}
		config.UnprotectedAuthRequests = append(config.UnprotectedAuthRequests, re)
	}

	config.FiberHealthCheckConfig = healthcheck.Config{
		LivenessEndpoint:  getEnvString("LIVENESS_ENDPOINT", "/live"),
		ReadinessEndpoint: getEnvString("READINESS_ENDPOINT", "/health"),
	}

	config.FiberRequestIDConfig = requestid.Config{
		ContextKey: getEnvString("REQUEST_ID_KEY", "requestid"),
	}

	config.RateLimitConfig = &RateLimitConfig{
		LoginRateLimitKey:    getEnvString("RATE_LIMIT_LOGIN_KEY", "login_attempts_rate_limit:"),
		LoginRateLimitTime:   getEnvDuration("RATE_LIMIT_LOGIN_TIME", 10*time.Second),
		LoginRateLimit:       getEnvInt("RATE_LIMIT_LOGIN_MAX_REQUESTS", 3),
		DefaultRateLimitKey:  getEnvString("RATE_LIMIT_DEFAULT_KEY", "rate_limit:"),
		DefaultRateLimitTime: getEnvDuration("RATE_LIMIT_DEFAULT_TIME", 1*time.Second),
		DefaultRateLimit:     getEnvInt("RATE_LIMIT_DEFAULT_MAX_REQUESTS", 4),
	}

	pathConfig.Other.Liveness = config.FiberHealthCheckConfig.LivenessEndpoint
	pathConfig.Other.Readliness = config.FiberHealthCheckConfig.ReadinessEndpoint

	config.Paths = pathConfig

	registerPathRegexp, _ := regexp.Compile(config.Paths.Authentication.Register)
	loginPathRegexp, _ := regexp.Compile(config.Paths.Authentication.Login)

	config.UnprotectedAuthRequests = append(
		config.UnprotectedAuthRequests,
		registerPathRegexp,
		loginPathRegexp,
	)

	if config.IsDebug {
		pprofPathRegexp, _ := regexp.Compile(config.Paths.Other.Pprof)
		config.UnprotectedAuthRequests = append(config.UnprotectedAuthRequests, pprofPathRegexp)
	}

	notInfoLogging = []string{
		pathConfig.Other.Liveness,
		pathConfig.Other.Readliness,
		pathConfig.Other.Pprof,
	}

	config.MetricsPort = getEnvInt("METRICS_PORT", 2112)

	config.SlowRequestThresholdMs = 300

	config.TracerConfig = &tracer.Config{
		Endpoint:           getEnvString("TRACER_ENDPOINT", "localhost:4317"),
		ServiceName:        getEnvString("TRACER_SERVICE_NAME", ""),
		ServiceVersion:     getEnvString("TRACER_SERVICE_VERSION", "v1"),
		ServiceEnvironment: config.EnvType,
	}

	return config
}

func getEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
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
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			labels[key] = value
		}
	}
	return labels
}
