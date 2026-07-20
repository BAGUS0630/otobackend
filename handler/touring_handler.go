package handler

import (
	"fmt"
	"otomeet-backend/config"
	"otomeet-backend/model"
	_ "otomeet-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// TouringRequest adalah format input untuk create dan update touring
type TouringRequest struct {
	NamaTouring string `json:"nama_touring" example:"Touring Lembang"`
	Tujuan      string `json:"tujuan" example:"Lembang, Bandung"`
	Tanggal     string `json:"tanggal" example:"2026-08-15"`
	Waktu       string `json:"waktu" example:"08:00"`
	Deskripsi   string `json:"deskripsi" example:"Touring santai akhir pekan"`
	Kuota       int    `json:"kuota" example:"50"`
	ImageURL    string `json:"image_url" example:"https://example.com/image.jpg"`
	LokasiAwal  string `json:"lokasi_awal" example:"Gasibu"`
	LokasiAkhir string `json:"lokasi_akhir" example:"Lembang Park"`
	HargaTiket  int    `json:"harga_tiket" example:"0"`
}

type TouringResponse struct {
	model.Touring
	PesertaCount int64 `json:"peserta_count"`
}

// GetAllTourings godoc
// @Summary      Lihat Semua Jadwal Touring
// @Description  Mengambil semua daftar agenda touring komunitas dengan pagination dan search (User & Admin)
// @Tags         Touring
// @Produce      json
// @Security     BearerAuth
// @Param        page     query  int     false  "Nomor halaman (default: 1)"
// @Param        limit    query  int     false  "Jumlah data per halaman (default: 10)"
// @Param        search   query  string  false  "Cari berdasarkan nama atau tujuan touring"
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      500      {object}  utils.SwaggerBasicResponse
// @Router       /api/touring [get]
func GetAllTourings(c *fiber.Ctx) error {
	// Default pagination values
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")
	sortField := c.Query("sort", "")
	sortOrder := c.Query("order", "asc")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var tourings []model.Touring
	var total int64

	// Query dengan filter search dan pagination
	query := config.DB
	if search != "" {
		query = query.Where("title ILIKE ? OR destination ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Hitung total data
	query.Model(&model.Touring{}).Count(&total)

	// Terapkan pengurutan
	if sortField != "" {
		orderStr := sortField
		if sortOrder == "desc" {
			orderStr += " DESC"
		} else {
			orderStr += " ASC"
		}
		query = query.Order(orderStr)
	} else {
		// Default: Akan datang diurutkan lebih dekat ke hari ini, yang sudah lewat di bawah
		query = query.Order("CASE WHEN departure_date >= to_char(CURRENT_DATE, 'YYYY-MM-DD') THEN 0 ELSE 1 END ASC")
		query = query.Order("CASE WHEN departure_date >= to_char(CURRENT_DATE, 'YYYY-MM-DD') THEN departure_date END ASC")
		query = query.Order("CASE WHEN departure_date < to_char(CURRENT_DATE, 'YYYY-MM-DD') THEN departure_date END DESC")
	}

	// Ambil data dengan limit dan offset
	if err := query.Limit(limit).Offset(offset).Find(&tourings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data touring"})
	}

	var response []TouringResponse
	for _, t := range tourings {
		var count int64
		config.DB.Model(&model.Registration{}).Where("touring_id = ?", t.ID).Count(&count)
		response = append(response, TouringResponse{
			Touring:      t,
			PesertaCount: count,
		})
	}

	return c.JSON(fiber.Map{
		"data":       response,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"total_page": (total + int64(limit) - 1) / int64(limit),
	})
}

// GetTouringByID godoc
// @Summary      Detail Jadwal Touring
// @Description  Mengambil detail satu agenda touring berdasarkan ID
// @Tags         Touring
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Touring ID"
// @Success      200  {object}  model.Touring
// @Failure      401  {object}  utils.Swagger401Response
// @Failure      404  {object}  utils.SwaggerBasicResponse
// @Router       /api/touring/{id} [get]
func GetTouringByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var touring model.Touring
	if err := config.DB.First(&touring, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Jadwal touring tidak ditemukan"})
	}
	return c.JSON(touring)
}

// CreateTouring godoc
// @Summary      Buat Agenda Touring Baru
// @Description  Menambahkan jadwal touring baru ke database (Khusus Admin Only)
// @Tags         Touring
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        touring  body      TouringRequest  true  "Data Agenda"
// @Success      201      {object}  utils.SwaggerBasicResponse
// @Failure      400      {object}  utils.SwaggerBasicResponse
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      403      {object}  utils.Swagger403Response
// @Failure      500      {object}  utils.SwaggerBasicResponse
// @Router       /api/touring [post]
func CreateTouring(c *fiber.Ctx) error {
	var input TouringRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// Validasi Backend[cite: 1]
	if input.NamaTouring == "" || input.Tujuan == "" || input.Kuota <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "Field wajib diisi dan kuota tidak boleh negatif"})
	}

	touring := model.Touring{
		NamaTouring: input.NamaTouring,
		Tujuan:      input.Tujuan,
		Tanggal:     input.Tanggal,
		Waktu:       input.Waktu,
		Deskripsi:   input.Deskripsi,
		Kuota:       input.Kuota,
		ImageURL:    input.ImageURL,
		LokasiAwal:  input.LokasiAwal,
		LokasiAkhir: input.LokasiAkhir,
		HargaTiket:  input.HargaTiket,
	}

	if err := config.DB.Create(&touring).Error; err != nil {
		fmt.Println("DB Create Error:", err)
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan jadwal touring"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Jadwal touring baru berhasil dibuat", "data": touring})
}

// UpdateTouring godoc
// @Summary      Edit Jadwal Touring
// @Description  Mengubah data agenda touring berdasarkan ID (Khusus Admin Only)
// @Tags         Touring
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int             true  "Touring ID"
// @Param        touring  body      TouringRequest  true  "Data Perubahan"
// @Success      200      {object}  utils.SwaggerBasicResponse
// @Failure      400      {object}  utils.SwaggerBasicResponse
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      403      {object}  utils.Swagger403Response
// @Failure      404      {object}  utils.SwaggerBasicResponse
// @Router       /api/touring/{id} [put]
func UpdateTouring(c *fiber.Ctx) error {
	id := c.Params("id")
	var touring model.Touring
	if err := config.DB.First(&touring, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Data tidak ditemukan"})
	}

	var input TouringRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// Update data menggunakan GORM
	config.DB.Model(&touring).Updates(model.Touring{
		NamaTouring: input.NamaTouring,
		Tujuan:      input.Tujuan,
		Tanggal:     input.Tanggal,
		Waktu:       input.Waktu,
		Deskripsi:   input.Deskripsi,
		Kuota:       input.Kuota,
		ImageURL:    input.ImageURL,
		LokasiAwal:  input.LokasiAwal,
		LokasiAkhir: input.LokasiAkhir,
		HargaTiket:  input.HargaTiket,
	})
	return c.JSON(fiber.Map{"message": "Jadwal touring berhasil diperbarui", "data": touring})
}

// DeleteTouring godoc
// @Summary      Hapus Jadwal Touring
// @Description  Menghapus agenda touring berdasarkan ID (Khusus Admin Only)
// @Tags         Touring
// @Security     BearerAuth
// @Param        id   path      int  true  "Touring ID"
// @Success      200  {object}  utils.SwaggerBasicResponse
// @Failure      401  {object}  utils.Swagger401Response
// @Failure      403  {object}  utils.Swagger403Response
// @Failure      404  {object}  utils.SwaggerBasicResponse
// @Router       /api/touring/{id} [delete]
func DeleteTouring(c *fiber.Ctx) error {
	id := c.Params("id")
	var touring model.Touring
	if err := config.DB.First(&touring, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Data tidak ditemukan"})
	}

	config.DB.Delete(&touring)
	return c.JSON(fiber.Map{"message": "Jadwal touring sukses dihapus"})
}

func JoinTouring(c *fiber.Ctx) error {
	// 1. Ambil ID touring dari URL parameter (:id)
	touringID := c.Params("id")

	// 2. Skenario Logika Sederhana (Mengurangi Kuota Langsung)
	// Catatan: Jika Anda punya tabel relasi pendaftaran, simpan datanya di sini.
	// Di bawah ini adalah contoh logika untuk mengurangi kuota langsung di Supabase lewat GORM:
	var touring model.Touring
	if err := config.DB.First(&touring, touringID).Error; err != nil {
		return c.Status(44).JSON(fiber.Map{"message": "Jadwal touring tidak ditemukan"})
	}

	if touring.Kuota <= 0 {
		return c.Status(400).JSON(fiber.Map{"message": "Maaf, kuota peserta untuk touring ini sudah penuh!"})
	}

	// Kurangi kuota 1 angka
	touring.Kuota = touring.Kuota - 1
	if err := config.DB.Save(&touring).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal bergabung ke touring"})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Berhasil bergabung ke agenda touring! 🏍️",
		"data":    touring,
	})
}
