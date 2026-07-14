package handler

import (
	"github.com/gofiber/fiber/v2"
	"otomeet-backend/config"
	"otomeet-backend/model"
)

// GetStats mengembalikan statistik singkat untuk dashboard frontend
func GetStats(c *fiber.Ctx) error {
	var usersCount int64
	var touringsCount int64
	var registrationsCount int64

	config.DB.Model(&model.User{}).Count(&usersCount)
	config.DB.Model(&model.Touring{}).Count(&touringsCount)
	config.DB.Model(&model.Registration{}).Count(&registrationsCount)

	return c.JSON(fiber.Map{
		"total_users":         usersCount,
		"total_tourings":      touringsCount,
		"total_registrations": registrationsCount,
	})
}
