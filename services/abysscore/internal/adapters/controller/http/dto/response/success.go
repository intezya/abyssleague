package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/intezya/abyssleague/services/abysscore/internal/domain/dto"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
	Path    string      `json:"path,omitempty"`
}

type PaginationResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Path    string `json:"path,omitempty"`

	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`

	Data interface{} `json:"data,omitempty"`
}

const successMessage = "success"

func Success(data interface{}, c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		Response{
			Message: successMessage,
			Data:    data,
			Code:    fiber.StatusOK,
			Path:    c.Path(),
		},
	)
}

func NoContent(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).JSON(
		Response{
			Message: successMessage,
			Code:    fiber.StatusNoContent,
			Path:    c.Path(),
		},
	)
}

func SuccessPagination[T any](data *dto.PaginatedResult[T], c *fiber.Ctx) error {
	// TODO: maybe set X-Total-Count
	return c.Status(fiber.StatusOK).JSON(
		PaginationResponse{
			Message: successMessage,
			Code:    fiber.StatusOK,
			Path:    c.Path(),

			Page:       data.Page,
			Size:       data.Size,
			TotalItems: data.TotalItems,
			TotalPages: data.TotalPages,
			Data:       data.Data,
		},
	)
}
