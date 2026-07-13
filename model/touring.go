package model

import "time"

type Touring struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	NamaTouring string    `gorm:"column:title;not null" json:"nama_touring"`
	Tujuan      string    `gorm:"column:destination;not null" json:"tujuan"`
	Tanggal     string    `gorm:"column:departure_date;not null" json:"tanggal"`
	Waktu       string    `json:"waktu"`
	Deskripsi   string    `gorm:"column:description;type:text" json:"deskripsi"`
	Kuota       int       `gorm:"column:max_participants" json:"kuota"`
	ImageURL    string    `gorm:"column:image_url" json:"image_url"`
	LokasiAwal  string    `json:"lokasi_awal"`
	LokasiAkhir string    `json:"lokasi_akhir"`
	HargaTiket  int       `json:"harga_tiket"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Registration struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id"`
	User         User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	TouringID    uint      `json:"touring_id"`
	BikeModel    string    `gorm:"not null" json:"bike_model"`
	LicensePlate string    `gorm:"not null" json:"license_plate"`
	RegisteredAt time.Time `gorm:"autoCreateTime" json:"registered_at"`
}

type Diskusi struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TouringID uint      `gorm:"not null" json:"touring_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
