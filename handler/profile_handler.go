package handler

import (
	"otomeet-backend/config"
	"otomeet-backend/model"
	"otomeet-backend/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// GetProfile godoc
// @Summary      Lihat Profil User
// @Description  Mengambil data profil user yang sedang login
// @Tags         User Profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  model.User
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/me [get]
func GetProfile(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Unauthorized", "UNAUTHORIZED")
	}
	userID := uint(userIDFloat)

	var user model.User
	if err := config.DB.Preload("Vehicles").First(&user, userID).Error; err != nil {
		return utils.RespondError(c, 404, "User tidak ditemukan", "USER_NOT_FOUND")
	}

	return utils.RespondSuccess(c, 200, "Berhasil mengambil profil", user)
}

// GetUserByID godoc
// @Summary      Get User By ID
// @Description  Mengambil data publik user berdasarkan ID
// @Tags         User Profile
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  utils.ErrorResponse
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
// @Param        profile  body  model.UpdateProfileRequest  true  "Data Profil"
// @Success      200  {object}  utils.SuccessResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Router       /api/profile [put]
func UpdateProfile(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Unauthorized", "UNAUTHORIZED")
	}
	userID := uint(userIDFloat)

	var input struct {
		FullName    string `json:"full_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}

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
// @Summary      Hapus Akun User
// @Description  Menghapus akun user yang sedang login setelah verifikasi password
// @Tags         User Profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      map[string]string  true  "Password konfirmasi"
// @Success      200   {object}  utils.SuccessResponse
// @Failure      401   {object}  utils.ErrorResponse
// @Router       /api/account [delete]
func DeleteAccount(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Unauthorized", "UNAUTHORIZED")
	}
	userID := uint(userIDFloat)

	var input struct {
		Password string `json:"password"`
	}

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
