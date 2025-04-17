package middleware

import (
	"abysscore/internal/adapters/config"
	adaptererror "abysscore/internal/common/errors/adapter"
	rediswrapper "abysscore/internal/infrastructure/cache/redis"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/pkglib/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"regexp"
	"time"
)

type RateLimitMiddleware struct {
	redisClient          *rediswrapper.ClientWrapper
	config               *config.Config
	authRouteRegexp      *regexp.Regexp
	loginAttemptsCounter *prometheus.CounterVec
	rateLimitCounter     *prometheus.CounterVec
}

// NewRateLimitMiddleware создает новый экземпляр RateLimitMiddleware
func NewRateLimitMiddleware(
	redisClient *rediswrapper.ClientWrapper,
	cfg *config.Config,
) *RateLimitMiddleware {
	loginAttemptsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_rate_limit_attempts",
			Help: "Count of login attempts hitting rate limits",
		},
		[]string{"username", "status"},
	)

	rateLimitCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_rate_limit_hits",
			Help: "Count of requests hitting general API rate limits",
		},
		[]string{"ip", "path", "status"},
	)

	prometheus.MustRegister(loginAttemptsCounter, rateLimitCounter)

	var authRouteRegexp *regexp.Regexp

	authRouteRegexp, err := regexp.Compile(cfg.Paths.Authentication.Self)

	if err != nil {
		logger.Log.Warn("Failed to get regexp from cfg.Paths.Authentication.Self: ", err)

		authRouteRegexp = nil
	}

	return &RateLimitMiddleware{
		redisClient:          redisClient,
		config:               cfg,
		loginAttemptsCounter: loginAttemptsCounter,
		rateLimitCounter:     rateLimitCounter,
		authRouteRegexp:      authRouteRegexp,
	}
}

func (r *RateLimitMiddleware) isAuthRoute(path string) bool {
	if r.authRouteRegexp != nil && r.authRouteRegexp.MatchString(path) {
		return true
	}
	return false
}

func (r *RateLimitMiddleware) HandleForAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if r.redisClient == nil {
			logger.Log.Warn("Redis client is not available in RateLimitMiddleware")
			return c.Next()
		}

		if !r.isAuthRoute(c.Path()) {
			return c.Next()
		}

		username := c.FormValue("username")
		if username == "" {
			type loginRequest struct {
				Username string `json:"username"`
			}

			var req loginRequest
			if err := c.BodyParser(&req); err != nil || req.Username == "" {
				return adaptererror.BadRequest.ToErrorResponse(c)
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
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

// HandleDefault возвращает middleware для общего ограничения частоты запросов по IP
func (r *RateLimitMiddleware) HandleDefault() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if r.redisClient == nil {
			logger.Log.Warn("Redis client is not available in RateLimitMiddleware")
			return c.Next()
		}

		ip := c.IP()
		if ip == "" {
			ip = "unknown"
		}

		path := c.Path()
		requestID := c.Locals(r.config.FiberRequestIDConfig.ContextKey)

		key := r.config.RateLimitConfig.DefaultRateLimitKey + ip
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		count, err := r.redisClient.Client.Get(ctx, key).Int()

		if err != nil && !errors.Is(err, redis.Nil) {
			logger.Log.With("error", err).Warn("Redis error in general rate limit middleware")
			return c.Next()
		}

		if count >= r.config.RateLimitConfig.DefaultRateLimit {
			r.rateLimitCounter.WithLabelValues(ip, path, "blocked").Inc()

			logger.Log.With(
				"ip", ip,
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

		r.rateLimitCounter.WithLabelValues(ip, path, "allowed").Inc()

		logger.Log.With(
			"ip", ip,
			"path", path,
			"count", count+1,
			"limit", r.config.RateLimitConfig.DefaultRateLimit,
			"request_id", requestID,
		).Debug("Request counted in rate limit")

		return c.Next()
	}
}
