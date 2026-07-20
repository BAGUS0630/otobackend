package handler

import (
	"fmt"
	"otomeet-backend/config"
	"otomeet-backend/model"
	"otomeet-backend/utils"

	"github.com/gofiber/fiber/v2"
)

// UploadProfilePhoto godoc
// @Summary      Upload Foto Profil
// @Description  Mengupload foto profil user ke Supabase Storage
// @Tags         User Profile
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file  formData  file  true  "File Foto Profil"
// @Success      200   {object}  utils.SuccessResponse
// @Failure      400   {object}  utils.ErrorResponse
// @Router       /api/upload-photo [post]
func UploadProfilePhoto(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Unauthorized", "UNAUTHORIZED")
	}
	userID := uint(userIDFloat)

	// Ambil file dari request
	file, err := c.FormFile("file")
	if err != nil {
		return utils.RespondError(c, 400, "File tidak ditemukan", "FILE_NOT_FOUND")
	}

	// Validasi dan upload file
	photoURL, err := utils.UploadToSupabase(file, "profile-photos", fmt.Sprintf("user-%d", userID))
	if err != nil {
		return utils.RespondError(c, 400, err.Error(), "UPLOAD_FAILED")
	}

	// Update profile photo di database
	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.RespondError(c, 404, "User tidak ditemukan", "USER_NOT_FOUND")
	}

	// Hapus foto lama jika ada
	if user.ProfilePhoto != "" {
		_ = utils.DeleteFromSupabase(user.ProfilePhoto, "profile-photos")
	}

	user.ProfilePhoto = photoURL
	if err := config.DB.Save(&user).Error; err != nil {
		return utils.RespondError(c, 500, "Gagal menyimpan foto profil", "SAVE_FAILED")
	}

	return utils.RespondSuccess(c, 200, "Foto profil berhasil diupload", fiber.Map{
		"profile_photo": photoURL,
	})
}

// DeleteProfilePhoto godoc
// @Summary      Hapus Foto Profil
// @Description  Menghapus foto profil user dari Supabase Storage
// @Tags         User Profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.SuccessResponse
// @Failure      400  {object}  utils.ErrorResponse
// @Router       /api/delete-photo [delete]
func DeleteProfilePhoto(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return utils.RespondError(c, 401, "Unauthorized", "UNAUTHORIZED")
	}
	userID := uint(userIDFloat)

	var user model.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.RespondError(c, 404, "User tidak ditemukan", "USER_NOT_FOUND")
	}

	if user.ProfilePhoto == "" {
		return utils.RespondError(c, 400, "Tidak ada foto profil untuk dihapus", "NO_PHOTO")
	}

	// Hapus dari Supabase
	if err := utils.DeleteFromSupabase(user.ProfilePhoto, "profile-photos"); err != nil {
		// Tetap lanjutkan, karena file mungkin sudah dihapus
	}

	user.ProfilePhoto = ""
	if err := config.DB.Save(&user).Error; err != nil {
		return utils.RespondError(c, 500, "Gagal menghapus foto profil", "DELETE_FAILED")
	}

	return utils.RespondSuccessNoData(c, 200, "Foto profil berhasil dihapus")
}
