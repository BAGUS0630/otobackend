package handler

import (
	"otomeet-backend/config"
	"otomeet-backend/model"
	"otomeet-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// VehicleRequest adalah format input untuk menambah kendaraan
type VehicleRequest struct {
	BikeModel    string `json:"bike_model" example:"Honda PCX"`
	LicensePlate string `json:"license_plate" example:"D 5678 EFG"`
}

// AddVehicle godoc
// @Summary      Tambah Kendaraan
// @Description  Menambahkan kendaraan baru ke garasi user
// @Tags         Vehicle
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        vehicle  body      VehicleRequest  true  "Data Kendaraan"
// @Success      201      {object}  utils.SwaggerBasicResponse
// @Failure      400      {object}  utils.SwaggerBasicResponse
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      500      {object}  utils.SwaggerBasicResponse
// @Router       /api/vehicles [post]
func AddVehicle(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Sesi login tidak valid", "INVALID_SESSION")
	}

	var input VehicleRequest

	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, 400, "Input tidak valid", "INVALID_INPUT")
	}

	vehicle := model.Vehicle{
		UserID:       uint(userIDFloat),
		BikeModel:    input.BikeModel,
		LicensePlate: input.LicensePlate,
	}

	if err := config.DB.Create(&vehicle).Error; err != nil {
		return utils.RespondError(c, 500, "Gagal menambahkan kendaraan", "CREATE_FAILED")
	}

	return utils.RespondSuccess(c, 201, "Kendaraan berhasil ditambahkan", vehicle)
}

// GetMyVehicles godoc
// @Summary      Lihat Garasi Kendaraan
// @Description  Mengambil daftar kendaraan milik user yang sedang login
// @Tags         Vehicle
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      500      {object}  utils.SwaggerBasicResponse
// @Router       /api/vehicles [get]
func GetMyVehicles(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Sesi login tidak valid", "INVALID_SESSION")
	}

	var vehicles []model.Vehicle
	if err := config.DB.Where("user_id = ?", uint(userIDFloat)).Find(&vehicles).Error; err != nil {
		return utils.RespondError(c, 500, "Gagal mengambil data kendaraan", "FETCH_FAILED")
	}

	// Sesuai dengan fetch di frontend: res.data.data
	return c.JSON(fiber.Map{"data": vehicles})
}

// DeleteVehicle godoc
// @Summary      Hapus Kendaraan
// @Description  Menghapus kendaraan dari garasi berdasarkan ID
// @Tags         Vehicle
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int  true  "Vehicle ID"
// @Success      200      {object}  utils.SwaggerBasicResponse
// @Failure      401      {object}  utils.Swagger401Response
// @Failure      404      {object}  utils.SwaggerBasicResponse
// @Failure      500      {object}  utils.SwaggerBasicResponse
// @Router       /api/vehicles/{id} [delete]
func DeleteVehicle(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Sesi login tidak valid", "INVALID_SESSION")
	}

	vehicleID := c.Params("id")

	var vehicle model.Vehicle
	if err := config.DB.Where("id = ? AND user_id = ?", vehicleID, uint(userIDFloat)).First(&vehicle).Error; err != nil {
		return utils.RespondError(c, 404, "Kendaraan tidak ditemukan", "NOT_FOUND")
	}

	if err := config.DB.Delete(&vehicle).Error; err != nil {
		return utils.RespondError(c, 500, "Gagal menghapus kendaraan", "DELETE_FAILED")
	}

	return utils.RespondSuccessNoData(c, 200, "Kendaraan berhasil dihapus")
}
