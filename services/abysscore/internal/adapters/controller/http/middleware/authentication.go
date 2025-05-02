package middleware

import (
	"context"
	"errors"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/apperrors"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
	domainservice "github.com/intezya/abyssleague/services/abysscore/internal/domain/service"
	rediswrapper "github.com/intezya/abyssleague/services/abysscore/internal/infrastructure/cache/redis"
	"github.com/intezya/pkglib/logger"
)

/* TODO:

1) 	Cache logic should use redis, not sync.Map (care about full user serializing,
	because most fields of UserDTO is necessary, but not all can be in http-response
	json and they are use tag json:"-")

*/

type CtxKey string

const authorizationHeaderKey = "Authorization"

const UserCtxKey CtxKey = "user"

const TokenCacheTime = 10 * time.Second

var (
	errTokenNotFoundInCache = errors.New("token not found in cache")
	errInvalidCacheType     = errors.New("invalid cache type")
)

type AuthenticationMiddleware struct {
	authenticationService domainservice.AuthenticationService
	redisClient           *rediswrapper.ClientWrapper

	localCache sync.Map
}

func NewAuthenticationMiddleware(
	authenticationService domainservice.AuthenticationService,
	redisClient *rediswrapper.ClientWrapper,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		authenticationService: authenticationService,
		redisClient:           redisClient,
		localCache:            sync.Map{},
	}
}

func (a *AuthenticationMiddleware) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Log.Debug("Starting authentication middleware")

		authorizationHeaderValue := c.Get(authorizationHeaderKey)
		authorizationHeaderValue = parseAuthorizationValue(authorizationHeaderValue)

		logger.Log.Debug("Authorization header value: ", authorizationHeaderValue)
		logger.Log.Debug("Redis client is: ", a.redisClient.Client)

		user, err := a.checkTokenCache(authorizationHeaderValue)
		if err != nil {
			logger.Log.Debug("Error checking token in cache: ", err)

			user, err = a.authenticationService.ValidateToken(
				c.UserContext(),
				authorizationHeaderValue,
			)
			if err != nil {
				logger.Log.Debug("Error validating token: ", err)

				return apperrors.HandleError(apperrors.WrapUnauthorized(err), c)
			}

			a.cacheToken(authorizationHeaderValue, user)
		}

		userContext := context.WithValue(c.UserContext(), UserCtxKey, user)
		c.SetUserContext(userContext)

		return c.Next()
	}
}

func (a *AuthenticationMiddleware) checkTokenCache(token string) (*dto.UserDTO, error) {
	val, ok := a.localCache.Load(token)

	if !ok {
		return nil, errTokenNotFoundInCache
	}

	user, ok := val.(*dto.UserDTO)

	if !ok {
		return nil, errInvalidCacheType
	}

	return user, nil
}

func (a *AuthenticationMiddleware) cacheToken(token string, user *dto.UserDTO) {
	a.localCache.Store(token, user)

	go func() {
		time.Sleep(TokenCacheTime)
		a.localCache.Delete(token)
	}()
}

func parseAuthorizationValue(authorizationHeaderValue string) string {
	switch {
	case strings.HasPrefix(authorizationHeaderValue, "Bearer "):
		return strings.TrimPrefix(authorizationHeaderValue, "Bearer ")
	default:
		return authorizationHeaderValue
	}
}
