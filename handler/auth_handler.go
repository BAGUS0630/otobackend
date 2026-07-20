package handler

import (
	"os"
	"otomeet-backend/config"
	"otomeet-backend/model"
	"otomeet-backend/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest adalah format input untuk registrasi
type RegisterRequest struct {
	Username string `json:"username" example:"budi123"`
	Password string `json:"password" example:"rahasia"`
	Role     string `json:"role" example:"user"`
}

// Register godoc
// @Summary      Registrasi Akun Baru
// @Description  Membuat akun user atau admin baru untuk platform OtoMeet
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      RegisterRequest  true  "Data Akun"
// @Success      201   {object}  utils.SwaggerBasicResponse
// @Failure      400   {object}  utils.SwaggerBasicResponse
// @Failure      500   {object}  utils.SwaggerBasicResponse
// @Router       /register [post]
func Register(c *fiber.Ctx) error {
	// 1. Buat struct penampung sementara khusus input JSON register
	var input RegisterRequest

	// 2. Baca data request body ke dalam struct temporary
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// Validasi input minimal
	if input.Username == "" || input.Password == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Username dan password wajib diisi"})
	}

	// Hash password dengan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal memproses password"})
	}

	// 3. Pindahkan data ke model.User asli sebelum disimpan ke database
	var user model.User
	user.Username = input.Username
	user.Password = string(hashedPassword)
	user.Role = input.Role
	user.Email = input.Username + "@otomeet.local" // Generate unique email so DB constraint won't fail

	// Jika di baris bawah kode Anda memakai input.Role, ganti menjadi user.Role, contoh:
	if user.Role == "" {
		user.Role = "user"
	}

	// --- Lanjutkan ke fungsi simpan repository Anda di bawah menggunakan objek `user`

	// Simpan ke database
	if err := repository.CreateUser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Username sudah terdaftar"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Registrasi anggota berhasil!"})
}

// Login godoc
// @Summary      Masuk Aplikasi (Login)
// @Description  Autentikasi akun untuk mendapatkan token JWT
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        login  body      model.LoginRequest  true  "Username dan Password"
// @Success      200    {object}  utils.SwaggerBasicResponse
// @Failure      400    {object}  utils.SwaggerBasicResponse
// @Failure      401    {object}  utils.Swagger401Response
// @Failure      500    {object}  utils.SwaggerBasicResponse
// @Router       /login [post]
func Login(c *fiber.Ctx) error {
	var input model.LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Input tidak valid"})
	}

	// Cari user berdasarkan username
	user, err := repository.GetUserByUsername(input.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Username atau password salah"})
	}

	// Cek kesesuaian hash password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Username atau password salah"})
	}

	// Membuat Token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Expired dalam 3 hari
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal membuat session login"})
	}

	return c.JSON(fiber.Map{
		"message": "Login sukses!",
		"token":   tokenString,
		"role":    user.Role,
	})
}

// ChangePassword godoc
// @Summary      Ubah Password User
// @Description  Mengubah password user yang sedang aktif login
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        password  body      model.PasswordChangeRequest  true  "Password Lama dan Baru"
// @Success      200       {object}  utils.SwaggerBasicResponse
// @Failure      400       {object}  utils.SwaggerBasicResponse
// @Failure      401       {object}  utils.Swagger401Response
// @Failure      404       {object}  utils.SwaggerBasicResponse
// @Failure      500       {object}  utils.SwaggerBasicResponse
// @Router       /api/change-password [put]
func ChangePassword(c *fiber.Ctx) error {
	// Ambil data user_id dari middleware JWT
	userIDFloat, ok := c.Locals("user_id").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"message": "Sesi login tidak valid"})
	}

	var input model.PasswordChangeRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Format input salah"})
	}

	if input.OldPassword == "" || input.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"message": "Password lama dan baru wajib diisi"})
	}

	// Ambil data user dari database
	var user model.User
	if err := config.DB.First(&user, uint(userIDFloat)).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	// Validasi kecocokan password lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		return c.Status(401).JSON(fiber.Map{"message": "Password lama Anda salah"})
	}

	// Hash password baru
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal memproses password baru"})
	}

	// Simpan perubahan ke Supabase
	config.DB.Model(&user).Update("password", string(newHashedPassword))

	return c.JSON(fiber.Map{"message": "Password Anda berhasil diperbarui!"})
}
