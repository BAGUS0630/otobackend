package handler

import (
	"otomeet-backend/config"
	"otomeet-backend/model"

	"github.com/gofiber/fiber/v2"
)

// GetDiskusiByTouring godoc
// @Summary      Ambil Diskusi Touring
// @Description  Mengambil semua komentar/diskusi untuk suatu agenda touring
// @Tags         Diskusi
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Touring ID"
// @Success      200  {array}   model.Diskusi
// @Router       /api/touring/{id}/diskusi [get]
func GetDiskusiByTouring(c *fiber.Ctx) error {
	touringID := c.Params("id")
	var diskusi []model.Diskusi

	// Preload "User" untuk mengambil data pengirim (username)
	if err := config.DB.Preload("User").Where("touring_id = ?", touringID).Order("created_at asc").Find(&diskusi).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal memuat diskusi"})
	}

	return c.JSON(diskusi)
}

// CreateDiskusi godoc
// @Summary      Kirim Komentar Diskusi
// @Description  Menambahkan komentar baru ke dalam agenda touring tertentu
// @Tags         Diskusi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int            true  "Touring ID"
// @Param        diskusi  body      model.Diskusi  true  "Pesan Diskusi"
// @Success      201      {object}  map[string]interface{}
// @Router       /api/touring/{id}/diskusi [post]
func CreateDiskusi(c *fiber.Ctx) error {
	touringID := c.Params("id")
	
	// Cek user dari token JWT
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid, silakan login ulang"})
	}

	var input struct {
		Message string `json:"message"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	if input.Message == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Pesan tidak boleh kosong"})
	}

	// Buat objek diskusi
	var touring model.Touring
	if err := config.DB.First(&touring, touringID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Agenda touring tidak ditemukan"})
	}

	diskusi := model.Diskusi{
		TouringID: touring.ID,
		UserID:    uint(userIDFloat),
		Message:   input.Message,
	}

	if err := config.DB.Create(&diskusi).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengirim pesan diskusi"})
	}

	// Preload User untuk di-return supaya frontend bisa langsung render namanya
	config.DB.Preload("User").First(&diskusi, diskusi.ID)

	return c.Status(201).JSON(fiber.Map{"message": "Pesan terkirim", "data": diskusi})
}
