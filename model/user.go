package model

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"unique;not null" json:"username"`
	Email         string    `gorm:"unique" json:"email"` // Changed: removed NOT NULL to handle existing data
	Password      string    `gorm:"not null" json:"-"`
	Role          string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	FullName      string    `gorm:"column:full_name" json:"full_name"`
	PhoneNumber   string    `json:"phone_number"`
	ProfilePhoto  string    `gorm:"column:profile_photo" json:"profile_photo"`
	Vehicles      []Vehicle `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"vehicles,omitempty"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
