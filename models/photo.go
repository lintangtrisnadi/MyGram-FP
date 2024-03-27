package models

import (
	"time"
)

// Photo merupakan model untuk foto yang diunggah oleh pengguna.
type Photo struct {
	ID        int64     `gorm:"primaryKey"`
	Title     string    `gorm:"size:100;not null"`  // Judul foto
	Caption   string    `gorm:"size:200"`           // Keterangan foto (opsional)
	PhotoURL  string    `gorm:"type:text;not null"` // URL gambar foto
	UserID    int64     `gorm:"not null"`           // ID pengguna yang mengunggah foto
	User      User      `gorm:"foreignKey:UserID"`  // Pengguna yang mengunggah foto
	CreatedAt time.Time // Waktu pembuatan foto
	UpdatedAt time.Time // Waktu pembaruan terakhir foto
	Comments  []Comment // Komentar yang terkait dengan foto
}

// CreatePhotoRequest adalah struktur data yang digunakan untuk membuat foto baru.
type CreatePhotoRequest struct {
	Title    string `json:"title" binding:"required"`     // Judul foto (wajib diisi)
	Caption  string `json:"caption,omitempty"`            // Keterangan foto (opsional)
	PhotoURL string `json:"photo_url" binding:"required"` // URL gambar foto (wajib diisi)
}

// UpdatePhoto adalah struktur data yang digunakan untuk memperbarui informasi foto yang sudah ada.
type UpdatePhoto struct {
	Title    string `json:"title" validate:"required"`     // Judul foto yang diperbarui (wajib diisi)
	Caption  string `json:"caption,omitempty"`             // Keterangan foto yang diperbarui (opsional)
	PhotoURL string `json:"photo_url" validate:"required"` // URL gambar foto yang diperbarui (wajib diisi)
}
