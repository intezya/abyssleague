package middleware

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/adapters/config"
	adaptererror "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/adapter"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/pkglib/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"time"
)

type RateLimitMiddleware struct {
	redisClient          *rediswrapper.ClientWrapper
	config               *config.Config
	loginAttemptsCounter *prometheus.CounterVec
	rateLimitCounter     *prometheus.CounterVec
}

const rateLimitStateReadTimeout = 2 * time.Second

func NewRateLimitMiddleware(
	redisClient *rediswrapper.ClientWrapper,
	cfg *config.Config,
) *RateLimitMiddleware {
	loginAttemptsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_rate_limit_attempts_total",
			Help: "Count of login attempts hitting rate limits",
		},
		[]string{"username", "status"},
	)

	rateLimitCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
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
		if r.redisClient == nil {
			logger.Log.Warn("Redis client is not available in RateLimitMiddleware")

			return c.Next()
		}

		username := c.FormValue("username")
		if username == "" {
			type loginRequest struct {
				Username string `json:"username"`
			}

			var req loginRequest
			if err := c.BodyParser(&req); err != nil || req.Username == "" {
				return c.Next() // will be handled as bad request
			}

			username = req.Username
		}

		requestID := c.Locals(r.config.FiberRequestIDConfig.ContextKey)

		logger.Log.With(
			"username", username,
			"path", c.Path(),
			"request_id", requestID,
		).Debug("Checking login rate limit")

		key := r.config.RateLimitConfig.LoginRateLimitKey + username

		ctx, cancel := context.WithTimeout(context.Background(), rateLimitStateReadTimeout)
		defer cancel()

		attempts, err := r.redisClient.Client.Get(ctx, key).Int()

		if errors.Is(err, redis.Nil) {
			err = r.redisClient.Client.SetEx(ctx, key, 1, r.config.RateLimitConfig.LoginRateLimitTime).Err()
			if err != nil {
				logger.Log.With("error", err).Warn("Failed to set rate limit counter")

				return c.Next()
			}

			r.loginAttemptsCounter.WithLabelValues(username, "allowed").Inc()

			logger.Log.With(
				"username", username,
				"attempts", 1,
				"limit", r.config.RateLimitConfig.LoginRateLimit,
				"request_id", requestID,
			).Debug("First login attempt recorded")

			return c.Next()
		}

		if err != nil {
			logger.Log.With("error", err).Warn("Redis error in rate limit middleware")

			return c.Next()
		}

		if attempts >= r.config.RateLimitConfig.LoginRateLimit {
			r.loginAttemptsCounter.WithLabelValues(username, "blocked").Inc()

			logger.Log.With(
				"username", username,
				"attempts", attempts,
				"limit", r.config.RateLimitConfig.LoginRateLimit,
				"request_id", requestID,
			).Info("Login rate limit exceeded")

			return adaptererror.TooManyRequests.ToErrorResponse(c)
		}

		_, err = r.redisClient.Client.Pipelined(
			ctx, func(pipe redis.Pipeliner) error {
				pipe.Incr(ctx, key)
				pipe.Expire(ctx, key, r.config.RateLimitConfig.LoginRateLimitTime)

				return nil
			},
		)

		if err != nil {
			logger.Log.With("error", err).Warn("Failed to increment rate limit counter")

			return c.Next()
		}

		r.loginAttemptsCounter.WithLabelValues(username, "allowed").Inc()

		logger.Log.With(
			"username", username,
			"attempts", attempts+1,
			"limit", r.config.RateLimitConfig.LoginRateLimit,
			"request_id", requestID,
		).Debug("Login attempt counted")

		return c.Next()
	}
}

// HandleDefault.
func (r *RateLimitMiddleware) HandleDefault() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if r.redisClient == nil || r.redisClient.Client == nil {
			logger.Log.Warn("Redis client is not available in RateLimitMiddleware")

			return c.Next()
		}

		ipAddr := c.IP()
		if ipAddr == "" {
			ipAddr = "unknown"
		}

		path := c.Path()
		requestID := c.Locals(r.config.FiberRequestIDConfig.ContextKey)

		key := r.config.RateLimitConfig.DefaultRateLimitKey + ipAddr

		ctx, cancel := context.WithTimeout(context.Background(), rateLimitStateReadTimeout)
		defer cancel()

		count, err := r.redisClient.Client.Get(ctx, key).Int()

		if err != nil && !errors.Is(err, redis.Nil) {
			logger.Log.With("error", err).Warn("Redis error in general rate limit middleware")

			return c.Next()
		}

		if count >= r.config.RateLimitConfig.DefaultRateLimit {
			r.rateLimitCounter.WithLabelValues(ipAddr, path, "blocked").Inc()

			logger.Log.With(
				"ip", ipAddr,
				"path", path,
				"count", count,
				"limit", r.config.RateLimitConfig.DefaultRateLimit,
				"request_id", requestID,
			).Warn("General rate limit exceeded")

			return adaptererror.TooManyRequests.ToErrorResponse(c)
		}

		_, err = r.redisClient.Client.Pipelined(
			ctx, func(pipe redis.Pipeliner) error {
				pipe.Incr(ctx, key)

				if count == 0 {
					pipe.Expire(ctx, key, r.config.RateLimitConfig.DefaultRateLimitTime)
				}

				return nil
			},
		)

		if err != nil {
			logger.Log.With("error", err).Warn("Failed to increment general rate limit counter")

			return c.Next()
		}

		r.rateLimitCounter.WithLabelValues(ipAddr, path, "allowed").Inc()

		logger.Log.With(
			"ipAddr", ipAddr,
			"path", path,
			"count", count+1,
			"limit", r.config.RateLimitConfig.DefaultRateLimit,
			"request_id", requestID,
		).Debug("Request counted in rate limit")

		return c.Next()
	}
}
