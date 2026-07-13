package handler

import (
	"github.com/gofiber/fiber/v2"
	"otomeet-backend/config"
)

// GetStats mengembalikan statistik singkat untuk dashboard frontend
func GetStats(c *fiber.Ctx) error {
	var usersCount int64
	var touringsCount int64
	var registrationsCount int64

	config.DB.Model(&map[string]interface{}{}).Table("users").Count(&usersCount)
	config.DB.Model(&map[string]interface{}{}).Table("tourings").Count(&touringsCount)
	config.DB.Model(&map[string]interface{}{}).Table("registrations").Count(&registrationsCount)

	return c.JSON(fiber.Map{
		"users":         usersCount,
		"tourings":      touringsCount,
		"registrations": registrationsCount,
	})
}
