package utils

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse adalah struktur standar untuk response error
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Code    string `json:"code,omitempty"`
}

// SuccessResponse adalah struktur standar untuk response sukses
type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondError mengirim response error yang terstandar
func RespondError(c *fiber.Ctx, statusCode int, message string, code string) error {
	return c.Status(statusCode).JSON(ErrorResponse{
		Status:  statusCode,
		Message: message,
		Code:    code,
	})
}

// RespondErrorWithDetail mengirim response error dengan detail error
func RespondErrorWithDetail(c *fiber.Ctx, statusCode int, message string, errDetail string, code string) error {
	return c.Status(statusCode).JSON(ErrorResponse{
		Status:  statusCode,
		Message: message,
		Error:   errDetail,
		Code:    code,
	})
}

// RespondSuccess mengirim response sukses
func RespondSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

// RespondSuccessNoData mengirim response sukses tanpa data
func RespondSuccessNoData(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(SuccessResponse{
		Status:  statusCode,
		Message: message,
	})
}
