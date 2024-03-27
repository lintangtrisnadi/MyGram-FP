package models

import (
	"time"
)

// User adalah model untuk pengguna dalam sistem.
type User struct {
	ID              int64         `gorm:"primaryKey"`         // ID pengguna
	Username        string        `gorm:"size:50;not null"`   // Nama pengguna
	Email           string        `gorm:"size:150;not null"`  // Email pengguna
	Password        string        `gorm:"type:text;not null"` // Kata sandi pengguna
	Age             int           `gorm:"not null"`           // Usia pengguna
	ProfileImageURL string        `gorm:"type:text"`          // URL gambar profil pengguna
	CreatedAt       time.Time     // Waktu pembuatan akun pengguna
	UpdatedAt       time.Time     // Waktu pembaruan terakhir akun pengguna
	Photos          []Photo       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Foto-foto yang dimiliki oleh pengguna
	Comments        []Comment     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Komentar yang dibuat oleh pengguna
	SocialMedias    []SocialMedia `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"` // Media sosial yang terkait dengan pengguna
}

// SignUpInput adalah struktur data yang digunakan saat mendaftar sebagai pengguna baru.
type SignUpInput struct {
	Username        string `json:"username" binding:"required"`            // Nama pengguna (wajib diisi)
	Email           string `json:"email" binding:"required"`               // Email pengguna (wajib diisi)
	Password        string `json:"password" binding:"required"`            // Kata sandi pengguna (wajib diisi)
	Age             int    `json:"age" binding:"required"`                 // Usia pengguna (wajib diisi)
	ProfileImageURL string `json:"profile_image_url" validate:"omitempty"` // URL gambar profil pengguna (opsional)
}

// SignInInput adalah struktur data yang digunakan saat masuk sebagai pengguna yang sudah terdaftar.
type SignInInput struct {
	Email    string `json:"email"  binding:"required"`    // Email pengguna yang sudah terdaftar (wajib diisi)
	Password string `json:"password"  binding:"required"` // Kata sandi pengguna yang sudah terdaftar (wajib diisi)
}

// UpdateCurrentUserRequest adalah struktur data yang digunakan untuk memperbarui informasi pengguna yang sudah masuk.
type UpdateCurrentUserRequest struct {
	Username        string `json:"username" binding:"required"`                      // Nama pengguna yang diperbarui (wajib diisi)
	Email           string `json:"email" binding:"required"`                         // Email pengguna yang diperbarui (wajib diisi)
	Age             int    `json:"age" binding:"required"`                           // Usia pengguna yang diperbarui (wajib diisi)
	ProfileImageURL string `json:"profile_image_url,omitempty" validate:"omitempty"` // URL gambar profil pengguna yang diperbarui (opsional)
}
