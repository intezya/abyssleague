package handlers

import (
	"abysscore/internal/adapters/controller/http/middleware"
	adaptererror "abysscore/internal/common/errors/adapter"
	"abysscore/internal/domain/dto"
	"context"
	"github.com/gofiber/fiber/v2"
)

func extractUserFromContext(ctx context.Context) (*dto.UserDTO, error) {
	user, ok := ctx.Value(middleware.UserCtxKey).(*dto.UserDTO)

	if !ok {
		return nil, adaptererror.InternalServerError
	}

	return user, nil
}

func extractIntParamOrErr(key string, c *fiber.Ctx) (int, error) {
	val, err := c.ParamsInt(key)

	if err != nil {
		return 0, adaptererror.BadRequestFunc(err)
	}

	return val, nil
}
