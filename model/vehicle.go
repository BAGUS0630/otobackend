package model

import "time"

type Vehicle struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	User         User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	BikeModel    string    `gorm:"column:bike_model;not null" json:"bike_model"`
	LicensePlate string    `gorm:"column:license_plate;not null" json:"license_plate"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
