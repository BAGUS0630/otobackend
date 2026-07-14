package config

import (
	"fmt"
	"otomeet-backend/model"

	"gorm.io/gorm"
)

// MigrateDB menangani migration dengan benar
func MigrateDB(db *gorm.DB) error {
	// 1. Auto migrate semua models dulu
	if err := db.AutoMigrate(
		&model.User{},
		&model.Touring{},
		&model.Registration{},
		&model.Diskusi{},
		&model.Vehicle{},
	); err != nil {
		return fmt.Errorf("gagal melakukan Auto Migrate: %w", err)
	}

	// 2. Update NULL email values setelah migration berhasil
	// Ini aman karena kolom email sudah dibuat (nullable atau tidak)
	db.Exec(`UPDATE "users" SET email = 'user_' || COALESCE(id, 0) || '@otomeet.local' WHERE email IS NULL OR email = ''`)

	return nil
}
