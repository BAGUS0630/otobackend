package handler

import (
	"otomeet-backend/config"
	"otomeet-backend/model"
	"otomeet-backend/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// UpdateProfileRequest adalah format input untuk update profil
type UpdateProfileRequest struct {
	FullName    string `json:"full_name" example:"Budi Santoso"`
	Email       string `json:"email" example:"budi@example.com"`
	PhoneNumber string `json:"phone_number" example:"08123456789"`
}

// DeleteAccountRequest adalah format input untuk hapus akun
type DeleteAccountRequest struct {
	Password string `json:"password" example:"rahasia123"`
}

// GetProfile godoc
// @Summary      Lihat Profil User
// @Description  Mengambil data profil user yang sedang login
// @Tags         User Profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  model.User
// @Failure      401  {object}  utils.Swagger401Response
// @Failure      404  {object}  utils.SwaggerBasicResponse
// @Failure      500  {object}  utils.SwaggerBasicResponse
// @Router       /api/me [get]
func GetProfile(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Sesi login tidak valid", "INVALID_SESSION")
	}
	userID := uint(userIDFloat)

	var user model.User
	if err := config.DB.Preload("Vehicles").First(&user, userID).Error; err != nil {
		return utils.RespondError(c, 404, "User tidak ditemukan", "USER_NOT_FOUND")
	}

	return utils.RespondSuccess(c, 200, "Berhasil mengambil profil", user)
}

// GetUserByID godoc
// @Summary      Lihat Profil User Lain
// @Description  Mengambil data profil user publik berdasarkan user ID
// @Tags         User Profile
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  utils.Swagger401Response
// @Failure      404  {object}  utils.SwaggerBasicResponse
// @Failure      500  {object}  utils.SwaggerBasicResponse
// @Router       /api/user/{id} [get]
func GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.RespondError(c, 404, "User tidak ditemukan", "USER_NOT_FOUND")
	}

	// Return hanya field public
	return c.JSON(fiber.Map{
		"id":            user.ID,
		"username":      user.Username,
		"full_name":     user.FullName,
		"profile_photo": user.ProfilePhoto,
		"created_at":    user.CreatedAt,
	})
}

// UpdateProfile godoc
// @Summary      Update Profil User
// @Description  Mengubah data profil user yang sedang login
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        profile  body  UpdateProfileRequest  true  "Data Profil"
// @Success      200  {object}  utils.SwaggerBasicResponse
// @Failure      400  {object}  utils.SwaggerBasicResponse
// @Failure      401  {object}  utils.Swagger401Response
// @Failure      404  {object}  utils.SwaggerBasicResponse
// @Failure      500  {object}  utils.SwaggerBasicResponse
// @Router       /api/profile [put]
func UpdateProfile(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Sesi login tidak valid", "INVALID_SESSION")
	}
	userID := uint(userIDFloat)

	var input UpdateProfileRequest

	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, 400, "Input tidak valid", "INVALID_INPUT")
	}

	// Validasi input
	if input.Email != "" && !utils.ValidateEmail(input.Email) {
		return utils.RespondError(c, 400, "Format email tidak valid", "INVALID_EMAIL")
	}

	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.RespondError(c, 404, "User tidak ditemukan", "USER_NOT_FOUND")
	}

	// Update field yang diisi
	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.PhoneNumber != "" {
		user.PhoneNumber = input.PhoneNumber
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return utils.RespondError(c, 400, "Gagal mengupdate profil", "UPDATE_FAILED")
	}

	return utils.RespondSuccess(c, 200, "Profil berhasil diperbarui", user)
}

// DeleteAccount godoc
// @Summary      Hapus Akun
// @Description  Menghapus akun user secara permanen
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        password  body  DeleteAccountRequest  true  "Password Konfirmasi"
// @Success      200  {object}  utils.SwaggerBasicResponse
// @Failure      400  {object}  utils.SwaggerBasicResponse
// @Failure      401  {object}  utils.Swagger401Response
// @Failure      404  {object}  utils.SwaggerBasicResponse
// @Failure      500  {object}  utils.SwaggerBasicResponse
// @Router       /api/delete-account [delete]
func DeleteAccount(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Sesi login tidak valid", "INVALID_SESSION")
	}
	userID := uint(userIDFloat)

	var input DeleteAccountRequest

	if err := c.BodyParser(&input); err != nil {
		return utils.RespondError(c, 400, "Input tidak valid", "INVALID_INPUT")
	}

	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.RespondError(c, 404, "User tidak ditemukan", "USER_NOT_FOUND")
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return utils.RespondError(c, 401, "Password tidak sesuai", "INVALID_PASSWORD")
	}

	// Hapus user
	if err := config.DB.Delete(&user).Error; err != nil {
		return utils.RespondError(c, 500, "Gagal menghapus akun", "DELETE_FAILED")
	}

	return utils.RespondSuccessNoData(c, 200, "Akun berhasil dihapus")
}
