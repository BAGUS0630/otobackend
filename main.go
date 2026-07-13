package main

import (
	"os"
	"otomeet-backend/config"
	"otomeet-backend/router"
	_ "otomeet-backend/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title OtoMeet API
// @version 1.0
// @description Ini adalah dokumentasi API untuk platform OtoMeet.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Koneksi ke Supabase
	config.ConnectDB()

	app := fiber.New()

	// 2. Pasang Middleware Wajib
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Nanti bisa diubah ke URL frontend setelah dideploy
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 3. Daftarkan Rute Endpoint
	router.SetupRoutes(app)

	// 4. Jalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	logFatal := app.Listen(port)
	if logFatal != nil {
		panic(logFatal)
	}
}