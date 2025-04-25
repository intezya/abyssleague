package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/config"
	adaptererror "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/adapter"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/pkglib/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

const rateLimitStateReadTimeout = 2 * time.Second

type RateLimitMiddleware struct {
	redisClient          *rediswrapper.ClientWrapper
	config               *config.Config
	loginAttemptsCounter *prometheus.CounterVec
	rateLimitCounter     *prometheus.CounterVec
}

func NewRateLimitMiddleware(
	redisClient *rediswrapper.ClientWrapper,
	cfg *config.Config,
) *RateLimitMiddleware {
	loginAttemptsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct // useless (for this application) fields
			Name: "login_rate_limit_attempts_total",
			Help: "Count of login attempts hitting rate limits",
		},
		[]string{"username", "status"},
	)

	rateLimitCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{ //nolint:exhaustruct // useless (for this application) fields
			Name: "api_rate_limit_hits_total",
			Help: "Count of requests hitting general API rate limits",
		},
		[]string{"ip", "path", "status"},
	)

	prometheus.MustRegister(loginAttemptsCounter, rateLimitCounter)

	return &RateLimitMiddleware{
		redisClient:          redisClient,
		config:               cfg,
		loginAttemptsCounter: loginAttemptsCounter,
		rateLimitCounter:     rateLimitCounter,
	}
}

func (r *RateLimitMiddleware) HandleForAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !r.isRedisAvailable() {
			return c.Next()
		}

		username := r.extractUsername(c)
		if username == "" {
			return c.Next() // will be handled as bad request
		}

		requestID := r.getRequestID(c)
		key := r.config.RateLimitConfig.LoginRateLimitKey + username

		return r.processLoginRateLimit(c, username, key, requestID)
	}
}

func (r *RateLimitMiddleware) HandleDefault() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !r.isRedisAvailable() {
			return c.Next()
		}

		ipAddr := r.getClientIP(c)
		path := c.Path()
		requestID := r.getRequestID(c)
		key := r.config.RateLimitConfig.DefaultRateLimitKey + ipAddr

		return r.processDefaultRateLimit(c, ipAddr, path, key, requestID)
	}
}

// Helper functions

func (r *RateLimitMiddleware) isRedisAvailable() bool {
	if r.redisClient == nil || r.redisClient.Client == nil {
		logger.Log.Warn("Redis client is not available in RateLimitMiddleware")

		return false
	}

	return true
}

func (r *RateLimitMiddleware) extractUsername(c *fiber.Ctx) string {
	username := c.FormValue("username")
	if username != "" {
		return username
	}

	// Try parsing from JSON body
	type loginRequest struct {
		Username string `json:"username"`
	}

	var req loginRequest
	if err := c.BodyParser(&req); err != nil || req.Username == "" {
		return ""
	}

	return req.Username
}

func (r *RateLimitMiddleware) getRequestID(c *fiber.Ctx) interface{} {
	return c.Locals(r.config.FiberRequestIDConfig.ContextKey)
}

func (r *RateLimitMiddleware) getClientIP(c *fiber.Ctx) string {
	ipAddr := c.IP()
	if ipAddr == "" {
		return "unknown"
	}

	return ipAddr
}

func (r *RateLimitMiddleware) processLoginRateLimit(
	c *fiber.Ctx,
	username, key string,
	requestID interface{},
) error {
	ctx, cancel := context.WithTimeout(context.Background(), rateLimitStateReadTimeout)
	defer cancel()

	attempts, err := r.redisClient.Client.Get(ctx, key).Int()

	// First login attempt
	if errors.Is(err, redis.Nil) {
		if err := r.setFirstLoginAttempt(ctx, key, username, requestID); err != nil {
			return c.Next()
		}

		return c.Next()
	}

	// Redis error
	if err != nil {
		logger.Log.With("error", err).Warn("Redis error in rate limit middleware")

		return c.Next()
	}

	// Check if limit exceeded
	if attempts >= r.config.RateLimitConfig.LoginRateLimit {
		r.logLimitExceeded(username, attempts, requestID)
		r.loginAttemptsCounter.WithLabelValues(username, "blocked").Inc()

		return adaptererror.TooManyRequests.ToErrorResponse(c)
	}

	// Increment attempt count
	if err := r.incrementLoginAttempt(ctx, key, username, attempts, requestID); err != nil {
		return c.Next()
	}

	return c.Next()
}

func (r *RateLimitMiddleware) setFirstLoginAttempt(
	ctx context.Context,
	key, username string,
	requestID interface{},
) error {
	err := r.redisClient.Client.SetEx(ctx, key, 1, r.config.RateLimitConfig.LoginRateLimitTime).
		Err()
	if err != nil {
		logger.Log.With("error", err).Warn("Failed to set rate limit counter")

		return err
	}

	r.loginAttemptsCounter.WithLabelValues(username, "allowed").Inc()

	logger.Log.With(
		"username", username,
		"attempts", 1,
		"limit", r.config.RateLimitConfig.LoginRateLimit,
		"request_id", requestID,
	).Debug("First login attempt recorded")

	return nil
}

func (r *RateLimitMiddleware) incrementLoginAttempt(
	ctx context.Context,
	key, username string,
	attempts int,
	requestID interface{},
) error {
	_, err := r.redisClient.Client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, r.config.RateLimitConfig.LoginRateLimitTime)

		return nil
	})
	if err != nil {
		logger.Log.With("error", err).Warn("Failed to increment rate limit counter")

		return err
	}

	r.loginAttemptsCounter.WithLabelValues(username, "allowed").Inc()

	logger.Log.With(
		"username", username,
		"attempts", attempts+1,
		"limit", r.config.RateLimitConfig.LoginRateLimit,
		"request_id", requestID,
	).Debug("Login attempt counted")

	return nil
}

func (r *RateLimitMiddleware) logLimitExceeded(
	username string,
	attempts int,
	requestID interface{},
) {
	logger.Log.With(
		"username", username,
		"attempts", attempts,
		"limit", r.config.RateLimitConfig.LoginRateLimit,
		"request_id", requestID,
	).Info("Login rate limit exceeded")
}

func (r *RateLimitMiddleware) processDefaultRateLimit(
	c *fiber.Ctx,
	ipAddr, path, key string,
	requestID interface{},
) error {
	ctx, cancel := context.WithTimeout(context.Background(), rateLimitStateReadTimeout)
	defer cancel()

	count, err := r.redisClient.Client.Get(ctx, key).Int()

	// Redis error
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Log.With("error", err).Warn("Redis error in general rate limit middleware")

		return c.Next()
	}

	// Check if limit exceeded
	if count >= r.config.RateLimitConfig.DefaultRateLimit {
		r.logDefaultLimitExceeded(ipAddr, path, count, requestID)
		r.rateLimitCounter.WithLabelValues(ipAddr, path, "blocked").Inc()

		return adaptererror.TooManyRequests.ToErrorResponse(c)
	}

	// Increment request count
	if err := r.incrementDefaultRequestCount(ctx, key, ipAddr, path, count, requestID); err != nil {
		return c.Next()
	}

	return c.Next()
}

func (r *RateLimitMiddleware) logDefaultLimitExceeded(
	ipAddr, path string,
	count int,
	requestID interface{},
) {
	logger.Log.With(
		"ip", ipAddr,
		"path", path,
		"count", count,
		"limit", r.config.RateLimitConfig.DefaultRateLimit,
		"request_id", requestID,
	).Warn("General rate limit exceeded")
}

func (r *RateLimitMiddleware) incrementDefaultRequestCount(
	ctx context.Context,
	key, ipAddr, path string,
	count int,
	requestID interface{},
) error {
	_, err := r.redisClient.Client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, key)

		if count == 0 {
			pipe.Expire(ctx, key, r.config.RateLimitConfig.DefaultRateLimitTime)
		}

		return nil
	})
	if err != nil {
		logger.Log.With("error", err).Warn("Failed to increment general rate limit counter")

		return err
	}

	r.rateLimitCounter.WithLabelValues(ipAddr, path, "allowed").Inc()

	logger.Log.With(
		"ip_address", ipAddr,
		"path", path,
		"count", count+1,
		"limit", r.config.RateLimitConfig.DefaultRateLimit,
		"request_id", requestID,
	).Debug("Request counted in rate limit")

	return nil
}
