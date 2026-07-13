package router

import (
	"otomeet-backend/handler"
	"otomeet-backend/middleware"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://127.0.0.1:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowCredentials: true,
	}))

	// Endpoint Dokumentasi Swagger UI (Publik)
	app.Get("/docs/*", swagger.HandlerDefault)

	// Rute Dasar Sistem
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to OtoMeet API! Server is running smoothly.",
			"status":  "Active",
		})
	})

	// Statistik singkat (publik)
	app.Get("/api/stats", handler.GetStats)

	// Endpoint Autentikasi Publik
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	// --- GROUP ROUTE TERPROTEKSI JWT ---
	api := app.Group("/api", middleware.Protected())

	// Fitur Akun User
	api.Get("/me", handler.GetProfile)                    // Lihat profil saya
	api.Put("/profile", handler.UpdateProfile)            // Update profil
	api.Put("/change-password", handler.ChangePassword)   // Ubah password
	api.Delete("/delete-account", handler.DeleteAccount)  // Hapus akun
	api.Get("/user/:id", handler.GetUserByID)             // Lihat profil user lain (public)
	api.Post("/upload-profile-photo", handler.UploadProfilePhoto)    // Upload foto profil
	api.Delete("/delete-profile-photo", handler.DeleteProfilePhoto)  // Hapus foto profil

	// Endpoint Fitur Utama 1 (Touring)
	api.Get("/touring", handler.GetAllTourings)     // Semua user login bisa melihat daftar agenda
	api.Get("/touring/:id", handler.GetTouringByID) // Detail data item touring

	// Fitur Diskusi / Forum Mini
	api.Get("/touring/:id/diskusi", handler.GetDiskusiByTouring)
	api.Post("/touring/:id/diskusi", handler.CreateDiskusi)

	api.Post("/touring/:id/join", handler.JoinTouring)

	// Endpoint Fitur Utama 2 (Registrasi / Gabung Komunitas)
	api.Post("/register-touring", handler.RegisterToTouring) // Mengikuti agenda touring
	api.Get("/my-touring", handler.GetMyTourings)            // Melihat riwayat touring saya

	// --- KHUSUS MANAJEMEN ROLE ADMIN ONLY ---
	adminRoutes := api.Group("", middleware.AdminOnly())
	adminRoutes.Post("/touring", handler.CreateTouring)              // Admin menambah jadwal
	adminRoutes.Put("/touring/:id", handler.UpdateTouring)           // Admin mengedit jadwal
	adminRoutes.Delete("/touring/:id", handler.DeleteTouring)        // Admin menghapus jadwal
	adminRoutes.Get("/touring/:id/peserta", handler.GetPesertaTouring) // Admin lihat daftar peserta
}