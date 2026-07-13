package handler

import (
	"otomeet-backend/config"
	"otomeet-backend/model"

	"github.com/gofiber/fiber/v2"
)

// RegisterToTouring godoc
// @Summary      Gabung Agenda Touring
// @Description  Mendaftarkan diri dan kendaraan ke dalam salah satu jadwal touring
// @Tags         Registration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        register  body      model.Registration  true  "Data Kendaraan"
// @Success      201       {object}  map[string]interface{}
// @Router       /api/register-touring [post]
func RegisterToTouring(c *fiber.Ctx) error {
	var registration model.Registration
	if err := c.BodyParser(&registration); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input data tidak valid"})
	}

	// Ambil userID otomatis dari JWT token yang sedang login
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid, silakan login ulang"})
	}
	registration.UserID = uint(userIDFloat)

	// Validasi input wajib
	if registration.TouringID == 0 || registration.BikeModel == "" || registration.LicensePlate == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Touring ID, model kendaraan, dan plat nomor wajib diisi"})
	}

	// Validasi: Cek apakah data agenda touring-nya eksis di database
	var touring model.Touring
	if err := config.DB.First(&touring, registration.TouringID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Agenda touring tidak ditemukan"})
	}

	// Validasi: Cek apakah user bersangkutan sudah mendaftar di touring yang sama
	var existingReg model.Registration
	errCheck := config.DB.Where("user_id = ? AND touring_id = ?", registration.UserID, registration.TouringID).First(&existingReg).Error
	if errCheck == nil {
		return c.Status(400).JSON(fiber.Map{"message": "Anda sudah terdaftar dalam agenda touring ini!"})
	}

	// Simpan pendaftaran ke database Supabase
	if err := config.DB.Create(&registration).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal melakukan registrasi touring"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Berhasil bergabung ke agenda touring!", "data": registration})
}

// GetMyTourings godoc
// @Summary      Riwayat Touring Saya
// @Description  Melihat daftar agenda touring apa saja yang pernah diikuti oleh user login saat ini
// @Tags         Registration
// @Produce      json
// @Security     BearerAuth
// @Success      200      {array}   model.Registration
// @Router       /api/my-touring [get]
func GetMyTourings(c *fiber.Ctx) error {
	userIDFloat := c.Locals("user_id").(float64)
	var registrations []model.Registration

	// Mengambil data dengan relasi (Preload) detail data Touring
	if err := config.DB.Preload("Touring").Where("user_id = ?", uint(userIDFloat)).Find(&registrations).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data"})
	}

	return c.JSON(registrations)
}

// GetPesertaTouring godoc
// @Summary      Daftar Peserta Touring
// @Description  Admin dapat melihat semua peserta yang mendaftar pada suatu agenda touring
// @Tags         Registration
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Touring ID"
// @Success      200  {array}   model.Registration
// @Router       /api/touring/{id}/peserta [get]
func GetPesertaTouring(c *fiber.Ctx) error {
	touringID := c.Params("id")
	var registrations []model.Registration

	if err := config.DB.Preload("User").Where("touring_id = ?", touringID).Find(&registrations).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data peserta"})
	}

	return c.JSON(registrations)
}
