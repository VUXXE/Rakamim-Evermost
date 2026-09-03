package helper

import (
	"github.com/gofiber/fiber/v2"
)

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginationData struct {
	Total  int64       `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Items  interface{} `json:"items"`
}

func Success(c *fiber.Ctx, code int, message string, data interface{}) error {
	return c.Status(code).JSON(BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Error(c *fiber.Ctx, code int, message string, errDetail interface{}) error {
	return c.Status(code).JSON(BaseResponse{
		Code:    code,
		Message: message,
		Data:    nil,
		Error:   errDetail,
	})
}

func Pagination(c *fiber.Ctx, code int, message string, total int64, limit, offset int, items interface{}) error {
	return c.Status(code).JSON(BaseResponse{
		Code:    code,
		Message: message,
		Data: PaginationData{
			Total:  total,
			Limit:  limit,
			Offset: offset,
			Items:  items,
		},
	})
}
