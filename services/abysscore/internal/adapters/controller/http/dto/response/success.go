package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
	Path    string      `json:"path,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(
		Response{
			Message: "success",
			Data:    data,
			Code:    fiber.StatusOK,
			Path:    c.Path(),
		},
	)
}
