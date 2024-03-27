package models

import (
	"time"
)

// SocialMedia adalah model untuk media sosial yang terkait dengan pengguna.
type SocialMedia struct {
	ID             int64     `gorm:"primaryKey"`         // ID media sosial
	Name           string    `gorm:"size:50;not null"`   // Nama media sosial
	SocialMediaURL string    `gorm:"type:text;not null"` // URL media sosial
	UserID         int64     `gorm:"not null"`           // ID pengguna yang terkait dengan media sosial
	User           User      `gorm:"foreignKey:UserID"`  // Pengguna yang terkait dengan media sosial
	CreatedAt      time.Time // Waktu pembuatan entri media sosial
	UpdatedAt      time.Time // Waktu pembaruan terakhir entri media sosial
}

// CreateSocialMediaRequest adalah struktur data yang digunakan untuk membuat entri media sosial baru.
type CreateSocialMediaRequest struct {
	Name           string `json:"name" validate:"required"`             // Nama media sosial (wajib diisi)
	SocialMediaURL string `json:"social_media_url" validate:"required"` // URL media sosial (wajib diisi)
}

// UpdateSocialMediaRequest adalah struktur data yang digunakan untuk memperbarui informasi entri media sosial yang sudah ada.
type UpdateSocialMediaRequest struct {
	Name           string `json:"name" validate:"required"`             // Nama media sosial yang diperbarui (wajib diisi)
	SocialMediaURL string `json:"social_media_url" validate:"required"` // URL media sosial yang diperbarui (wajib diisi)
}
