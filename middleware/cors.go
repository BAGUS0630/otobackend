package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Cors() fiber.Handler {
	return func(c *fiber.Ctx) error {
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:5173,http://127.0.0.1:5173"
		}

		origin := c.Get("Origin")
		allowedList := strings.Split(allowedOrigins, ",")

		for _, allowed := range allowedList {
			if strings.TrimSpace(allowed) == origin {
				c.Set("Access-Control-Allow-Origin", origin)
				c.Set("Access-Control-Allow-Credentials", "true")
				c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
				c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				break
			}
		}

		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
