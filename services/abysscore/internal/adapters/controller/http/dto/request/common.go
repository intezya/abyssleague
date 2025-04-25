package request

import (
	"github.com/gofiber/fiber/v2"
	adaptererror "github.com/intezya/abyssleague/services/abysscore/internal/common/errors/adapter"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/entity/types"
	"github.com/intezya/abyssleague/services/abysscore/internal/pkg/queryparser"
)

const (
	defaultPage = 1
	defaultSize = 10
)

type OrderTypeValidator[T ~string] = func(v string) (T, error)

type PaginationQuery[orderByT ~string] struct {
	Page      int
	Size      int
	OrderBy   orderByT
	OrderType types.OrderType
}

func NewPaginationQuery[orderByT ~string](
	c *fiber.Ctx,
	orderTypeValidator OrderTypeValidator[orderByT],
) (*PaginationQuery[orderByT], error) {
	orderBy, err := orderTypeValidator(c.Query("order_by", ""))
	if err != nil {
		return nil, adaptererror.BadRequestFunc(err)
	}

	orderType, err := queryparser.ParseOrderType(c.Query("order_type", ""))
	if err != nil {
		return nil, adaptererror.BadRequestFunc(err)
	}

	return &PaginationQuery[orderByT]{
		Page:      c.QueryInt("page", defaultPage),
		Size:      c.QueryInt("size", defaultSize),
		OrderBy:   orderBy,
		OrderType: orderType,
	}, nil
}
