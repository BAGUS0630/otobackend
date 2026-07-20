package handler

import (
	"otomeet-backend/config"
	"otomeet-backend/model"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "otomeet-backend/utils"
)

// RegisterTouringRequest adalah format input pendaftaran touring
type RegisterTouringRequest struct {
	TouringID    uint   `json:"touring_id" example:"1"`
	BikeModel    string `json:"bike_model" example:"Yamaha NMAX 155"`
	LicensePlate string `json:"license_plate" example:"D 1234 ABC"`
}

// RegisterToTouring godoc
// @Summary      Gabung Agenda Touring
// @Description  Mendaftarkan diri dan kendaraan ke dalam salah satu jadwal touring
// @Tags         Registration
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        register  body      RegisterTouringRequest  true  "Data Kendaraan"
// @Success      201       {object}  utils.SwaggerBasicResponse
// @Failure      400       {object}  utils.SwaggerBasicResponse
// @Failure      401       {object}  utils.Swagger401Response
// @Failure      404       {object}  utils.SwaggerBasicResponse
// @Failure      500       {object}  utils.SwaggerBasicResponse
// @Router       /api/register-touring [post]
func RegisterToTouring(c *fiber.Ctx) error {
	var input RegisterTouringRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input data tidak valid"})
	}

	// Ambil userID otomatis dari JWT token yang sedang login
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid, silakan login ulang"})
	}
	
	// Validasi input wajib
	if input.TouringID == 0 || input.BikeModel == "" || input.LicensePlate == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Touring ID, model kendaraan, dan plat nomor wajib diisi"})
	}

	// Validasi: Cek apakah data agenda touring-nya eksis di database
	var touring model.Touring
	if err := config.DB.First(&touring, input.TouringID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Agenda touring tidak ditemukan"})
	}

	// Validasi: Cek apakah touring sudah selesai
	if len(touring.Tanggal) >= 10 {
		tglStr := touring.Tanggal[:10]
		if tgl, errDate := time.Parse("2006-01-02", tglStr); errDate == nil {
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			if tgl.Before(today) {
				return c.Status(400).JSON(fiber.Map{"message": "Maaf, agenda touring ini sudah selesai dan tidak bisa diikuti lagi!"})
			}
		}
	}

	// Validasi: Cek apakah kuota penuh
	if touring.Kuota > 0 {
		var count int64
		config.DB.Model(&model.Registration{}).Where("touring_id = ?", touring.ID).Count(&count)
		if count >= int64(touring.Kuota) {
			return c.Status(400).JSON(fiber.Map{"message": "Maaf, kuota peserta untuk touring ini sudah penuh!"})
		}
	}

	// Validasi: Cek apakah user bersangkutan sudah mendaftar di touring yang sama
	var existingReg model.Registration
	errCheck := config.DB.Where("user_id = ? AND touring_id = ?", uint(userIDFloat), input.TouringID).First(&existingReg).Error
	if errCheck == nil {
		return c.Status(400).JSON(fiber.Map{"message": "Anda sudah terdaftar dalam agenda touring ini!"})
	}

	registration := model.Registration{
		UserID:       uint(userIDFloat),
		TouringID:    input.TouringID,
		BikeModel:    input.BikeModel,
		LicensePlate: input.LicensePlate,
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
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      500      {object}  utils.SwaggerBasicResponse
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
// @Failure      401  {object}  utils.Swagger401Response
// @Failure      403  {object}  utils.Swagger403Response
// @Failure      500  {object}  utils.SwaggerBasicResponse
// @Router       /api/touring/{id}/peserta [get]
func GetPesertaTouring(c *fiber.Ctx) error {
	touringID := c.Params("id")
	var registrations []model.Registration

	if err := config.DB.Preload("User").Where("touring_id = ?", touringID).Find(&registrations).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data peserta"})
	}

	return c.JSON(registrations)
}
