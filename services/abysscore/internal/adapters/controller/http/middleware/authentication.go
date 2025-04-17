package middleware

import (
	adaptererror "abysscore/internal/common/errors/adapter"
	"abysscore/internal/domain/dto"
	"abysscore/internal/domain/service"
	rediswrapper "abysscore/internal/infrastructure/cache/redis"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/pkglib/logger"
	"regexp"
	"strings"
	"time"
)

type ctxKey string

const authorizationHeaderKey = "Authorization"

const UserCtxKey ctxKey = "user"

type AuthenticationMiddleware struct {
	unprotectedRoutes     []*regexp.Regexp
	authenticationService domainservice.AuthenticationService
	redisClient           *rediswrapper.ClientWrapper
}

func NewAuthenticationMiddleware(
	unprotectedRoutes []*regexp.Regexp,
	authenticationService domainservice.AuthenticationService,
	redisClient *rediswrapper.ClientWrapper,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		unprotectedRoutes:     unprotectedRoutes,
		authenticationService: authenticationService,
		redisClient:           redisClient,
	}
}

func (a *AuthenticationMiddleware) isUnprotectedRoute(path string) bool {
	for _, re := range a.unprotectedRoutes {
		if re != nil && re.MatchString(path) {
			return true
		}
	}
	return false
}

func (a *AuthenticationMiddleware) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Log.Debug("Starting authentication middleware")

		if a.isUnprotectedRoute(c.Path()) {
			return c.Next()
		}

		authorizationHeaderValue := c.Get(authorizationHeaderKey)
		authorizationHeaderValue = parseAuthorizationValue(authorizationHeaderValue)

		logger.Log.Debug("Authorization header value: ", authorizationHeaderValue)
		logger.Log.Debug("Redis client is: ", a.redisClient.Client)

		user, err := a.checkTokenCache(authorizationHeaderValue)

		if err != nil {
			logger.Log.Debug("Error checking token in cache: ", err)
			user, err = a.authenticationService.ValidateToken(c.Context(), authorizationHeaderValue)

			if err != nil {
				logger.Log.Debug("Error validating token: ", err)
				return adaptererror.ErrUnauthorized(err).ToErrorResponse(c)
			}

			a.cacheToken(authorizationHeaderValue, user)
		}

		userContext := context.WithValue(c.UserContext(), UserCtxKey, user)
		c.SetUserContext(userContext)

		return c.Next()
	}
}

func (a *AuthenticationMiddleware) checkTokenCache(token string) (*dto.UserDTO, error) {
	ctx := context.Background()

	cachedUser, err := a.redisClient.Client.Get(ctx, token).Result()

	if err == nil {
		return deserializeUser(cachedUser), nil
	}

	if !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("redis error: %v", err)
	}

	return nil, nil // not in cache
}

func (a *AuthenticationMiddleware) cacheToken(token string, user *dto.UserDTO) {
	ctx := context.Background()

	serializedUser := serializeUser(user)

	err := a.redisClient.Client.Set(ctx, token, serializedUser, 10*time.Minute).Err()

	if err != nil {
		logger.Log.Error("Error caching token: ", err)
	}
}

func serializeUser(user *dto.UserDTO) string {
	b, _ := json.Marshal(user)
	return string(b)
}

func deserializeUser(data string) *dto.UserDTO {
	var user dto.UserDTO
	_ = json.Unmarshal([]byte(data), &user)
	return &user
}

func parseAuthorizationValue(authorizationHeaderValue string) string {
	switch {
	case strings.HasPrefix(authorizationHeaderValue, "Bearer "):
		return strings.TrimPrefix(authorizationHeaderValue, "Bearer ")
	default:
		return authorizationHeaderValue
	}
}
