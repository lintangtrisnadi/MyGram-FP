package models

import (
	"time"
)

// Comment merupakan model untuk komentar pada foto.
type Comment struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	PhotoID   int64     `gorm:"not null"`
	Photo     Photo     `gorm:"foreignKey:PhotoID"`
	Message   string    `gorm:"size:200;not null"`
	CreatedAt time.Time // Waktu pembuatan komentar
	UpdatedAt time.Time // Waktu pembaruan terakhir komentar
}

// CreateCommentRequest adalah struktur data yang digunakan untuk membuat komentar baru.
type CreateCommentRequest struct {
	PhotoID int64  `json:"photo_id" validate:"required"` // ID foto yang akan dikomentari
	Message string `json:"message" validate:"required"`  // Isi pesan komentar
}

// UpdateCommentRequest adalah struktur data yang digunakan untuk memperbarui komentar yang sudah ada.
type UpdateCommentRequest struct {
	Message string `json:"message" validate:"required"` // Isi pesan komentar yang diperbarui
}
