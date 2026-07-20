package handler

import (
	"otomeet-backend/config"
	"otomeet-backend/model"

	"github.com/gofiber/fiber/v2"
	_ "otomeet-backend/utils"
)

// DiskusiRequest adalah format input untuk pesan diskusi
type DiskusiRequest struct {
	Message string `json:"message" example:"Halo semuanya, kapan kita kumpul?"`
}

// GetDiskusiByTouring godoc
// @Summary      Ambil Diskusi Touring
// @Description  Mengambil semua komentar/diskusi untuk suatu agenda touring
// @Tags         Diskusi
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Touring ID"
// @Success      200  {array}   model.Diskusi
// @Failure      401  {object}  utils.Swagger401Response
// @Failure      500  {object}  utils.SwaggerBasicResponse
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
// @Param        id       path      int             true  "Touring ID"
// @Param        diskusi  body      DiskusiRequest  true  "Pesan Diskusi"
// @Success      201      {object}  utils.SwaggerBasicResponse
// @Failure      400      {object}  utils.SwaggerBasicResponse
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      404      {object}  utils.SwaggerBasicResponse
// @Failure      500      {object}  utils.SwaggerBasicResponse
// @Router       /api/touring/{id}/diskusi [post]
func CreateDiskusi(c *fiber.Ctx) error {
	touringID := c.Params("id")

	// Cek user dari token JWT
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid, silakan login ulang"})
	}

	var input DiskusiRequest

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

// GetLatestDiskusi godoc
// @Summary      Ambil Pesan Diskusi Terbaru Global
// @Description  Mengambil semua pesan terbaru di atas ID tertentu
// @Tags         Diskusi
// @Produce      json
// @Security     BearerAuth
// @Param        last_id query int true "Last Message ID"
// @Success      200 {array} model.Diskusi
// @Failure      401 {object} utils.Swagger401Response
// @Failure      500 {object} utils.SwaggerBasicResponse
// @Router       /api/diskusi/latest [get]
func GetLatestDiskusi(c *fiber.Ctx) error {
	lastID := c.Query("last_id", "0")
	var diskusi []model.Diskusi

	// Preload User dan Touring agar frontend tahu pesan ini dari grup mana
	if err := config.DB.Preload("User").Preload("Touring").
		Where("id > ?", lastID).
		Order("id asc").
		Find(&diskusi).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal memuat pesan terbaru"})
	}

	return c.JSON(diskusi)
}

// UpdateDiskusi godoc
// @Summary      Edit Pesan Diskusi
// @Description  Memperbarui teks pesan diskusi yang sudah ada. Hanya pemilik pesan yang bisa mengeditnya.
// @Tags         Diskusi
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int             true  "ID Diskusi"
// @Param        diskusi  body      DiskusiRequest  true  "Teks Pesan Baru"
// @Success      200      {object}  utils.SwaggerBasicResponse
// @Failure      400      {object}  utils.SwaggerBasicResponse
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      403      {object}  utils.Swagger403Response
// @Failure      404      {object}  utils.SwaggerBasicResponse
// @Failure      500      {object}  utils.SwaggerBasicResponse
// @Router       /api/diskusi/{id} [put]
func UpdateDiskusi(c *fiber.Ctx) error {
	diskusiID := c.Params("id")
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid"})
	}
	userID := uint(userIDFloat)

	var input DiskusiRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	var diskusi model.Diskusi
	if err := config.DB.First(&diskusi, diskusiID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Pesan tidak ditemukan"})
	}

	// Validasi kepemilikan pesan
	if diskusi.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak: Hanya pengirim yang bisa mengedit pesan ini"})
	}

	// Update pesan
	diskusi.Message = input.Message
	diskusi.IsEdited = true
	if err := config.DB.Save(&diskusi).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan perubahan pesan"})
	}

	return c.JSON(fiber.Map{"message": "Pesan berhasil diedit", "data": diskusi})
}

// DeleteDiskusi godoc
// @Summary      Hapus Pesan Diskusi
// @Description  Menghapus pesan diskusi. Hanya pemilik pesan yang bisa menghapusnya.
// @Tags         Diskusi
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      int  true  "ID Diskusi"
// @Success      200 {object}  utils.SwaggerBasicResponse
// @Failure      401 {object}  utils.Swagger401Response
// @Failure      403 {object}  utils.Swagger403Response
// @Failure      404 {object}  utils.SwaggerBasicResponse
// @Failure      500 {object}  utils.SwaggerBasicResponse
// @Router       /api/diskusi/{id} [delete]
func DeleteDiskusi(c *fiber.Ctx) error {
	diskusiID := c.Params("id")
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi tidak valid"})
	}
	userID := uint(userIDFloat)

	var diskusi model.Diskusi
	if err := config.DB.First(&diskusi, diskusiID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Pesan tidak ditemukan"})
	}

	// Validasi kepemilikan pesan
	if diskusi.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"message": "Akses ditolak: Hanya pengirim yang bisa menghapus pesan ini"})
	}

	if err := config.DB.Delete(&diskusi).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menghapus pesan"})
	}

	return c.JSON(fiber.Map{"message": "Pesan berhasil dihapus"})
}
