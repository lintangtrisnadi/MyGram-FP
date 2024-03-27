package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mygram-final-project/models"
	"mygram-final-project/utils"
)

// UserController mengelola operasi terkait pengguna
type UserController struct {
	DB *gorm.DB
}

// NewUserController membuat instance baru dari UserController
func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

// UpdateMe mengupdate informasi pengguna saat ini
func (uc *UserController) UpdateMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload models.UpdateCurrentUserRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Validasi email
	if !utils.IsValidEmail(payload.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format emai tidak valid."})
		return
	}

	// Validasi username
	if payload.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username harus diisi."})
		return
	}

	// Validasi usia
	if payload.Age < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Age minimal 8 tahun."})
		return
	}

	// Validasi URL gambar profil
	if payload.ProfileImageURL != "" && !utils.IsValidURL(payload.ProfileImageURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format profile image URL tidak valid."})
		return
	}

	// Memperbarui informasi pengguna saat ini
	currentUser.Username = payload.Username
	currentUser.Email = payload.Email
	currentUser.Age = payload.Age
	currentUser.ProfileImageURL = payload.ProfileImageURL

	// Validasi Username yang Unik
	existingUser := models.User{}
	if err := uc.DB.Where("username = ?", payload.Username).First(&existingUser).Error; err == nil && existingUser.ID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username sudah terdaftar."})
		return
	}

	// Validasi Email yang Unik
	if err := uc.DB.Where("email = ?", payload.Email).First(&existingUser).Error; err == nil && existingUser.ID != currentUser.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Email sudah terdaftar."})
		return
	}

	// Simpan perubahan ke database
	if err := uc.DB.Save(&currentUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengupdate informasi user."})
		return
	}

	// Siapkan respons sesuai spesifikasi OpenAPI
	responseData := gin.H{
		"id":                currentUser.ID,
		"email":             currentUser.Email,
		"username":          currentUser.Username,
		"age":               currentUser.Age,
		"profile_image_url": currentUser.ProfileImageURL,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// DeleteMe menghapus pengguna saat ini
func (uc *UserController) DeleteMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// Hapus pengguna saat ini dari database
	if err := uc.DB.Delete(&currentUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghapus user."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
