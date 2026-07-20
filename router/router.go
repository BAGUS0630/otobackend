package router

import (
	"otomeet-backend/handler"
	"otomeet-backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(middleware.Cors())

	// Endpoint Dokumentasi Swagger UI (Publik)
	app.Get("/docs/*", swagger.HandlerDefault)

	// Rute Dasar Sistem
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to OtoMeet API! Server is running smoothly.",
			"status":  "Active",
		})
	})

	// Endpoint Autentikasi Publik
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	// --- GROUP ROUTE TERPROTEKSI JWT ---
	api := app.Group("/api", middleware.Protected())

	// Fitur Akun User
	api.Get("/stats", handler.GetDashboardStats)          // Statistik Dashboard
	api.Get("/me", handler.GetProfile)                    // Lihat profil saya
	api.Get("/user/:id", handler.GetUserByID)             // Lihat profil member lain
	api.Put("/profile", handler.UpdateProfile)            // Update profil
	api.Put("/change-password", handler.ChangePassword)   // Ubah password

	// Fitur Kendaraan (Garasi)
	api.Get("/vehicles", handler.GetMyVehicles)
	api.Post("/vehicles", handler.CreateVehicle)
	api.Delete("/vehicles/:id", handler.DeleteVehicle)

	// Endpoint Fitur Utama 1 (Touring)
	api.Get("/touring", handler.GetAllTourings)     // Semua user login bisa melihat daftar agenda
	api.Get("/touring/:id", handler.GetTouringByID) // Detail data item touring

	// Fitur Diskusi / Forum Mini
	api.Get("/touring/:id/diskusi", handler.GetDiskusiByTouring)
	api.Post("/touring/:id/diskusi", handler.CreateDiskusi)
	api.Get("/diskusi/:id", handler.GetDiskusiByID)
	api.Put("/diskusi/:id", handler.UpdateDiskusi)
	api.Delete("/diskusi/:id", handler.DeleteDiskusi)


	// Endpoint Fitur Utama 2 (Registrasi / Gabung Komunitas)
	api.Post("/register-touring", handler.RegisterToTouring) // Mengikuti agenda touring
	api.Get("/my-touring", handler.GetMyTourings)            // Melihat riwayat touring saya

	// --- KHUSUS MANAJEMEN ROLE ADMIN ONLY ---
	adminRoutes := api.Group("", middleware.AdminOnly())
	adminRoutes.Post("/touring", handler.CreateTouring)              // Admin menambah jadwal
	adminRoutes.Put("/touring/:id", handler.UpdateTouring)           // Admin mengedit jadwal
	adminRoutes.Delete("/touring/:id", handler.DeleteTouring)        // Admin menghapus jadwal
	adminRoutes.Get("/touring/:id/peserta", handler.GetPesertaTouring) // Admin lihat daftar peserta
	adminRoutes.Delete("/touring/:id/peserta/:peserta_id", handler.DeletePesertaTouring) // Admin keluarkan peserta
}