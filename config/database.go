package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file tidak ditemukan, membaca system env")
	}

	dsn := os.Getenv("SUPABASE_DSN")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disables implicit prepared statement usage, needed for Supabase pooler
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database Supabase:", err)
	}

	fmt.Println("🚀 Berhasil terkoneksi ke PostgreSQL Supabase!")

	// Auto Migration untuk membuat tabel otomatis dengan handling proper
	err = MigrateDB(db)
	if err != nil {
		log.Fatal("Gagal melakukan Auto Migrate:", err)
	}

	DB = db
}
