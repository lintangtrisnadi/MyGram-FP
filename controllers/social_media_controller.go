package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mygram-final-project/models"
	"mygram-final-project/utils"
)

// SocialMediaController struct digunakan untuk mengendalikan operasi-operasi terkait media sosial
type SocialMediaController struct {
	DB *gorm.DB
}

// NewSocialMediaController membuat instance baru dari SocialMediaController
func NewSocialMediaController(DB *gorm.DB) SocialMediaController {
	return SocialMediaController{DB}
}

// CreateSocialMedia membuat entri baru untuk media sosial
func (smc *SocialMediaController) CreateSocialMedia(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload models.CreateSocialMediaRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validasi field name
	if payload.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Name harus diisi."})
		return
	}

	// Validasi field social_media_url
	if payload.SocialMediaURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Social media URL harus diisi."})
		return
	}

	// Validasi URL profil
	if payload.SocialMediaURL != "" && !utils.IsValidURL(payload.SocialMediaURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format profile image URL tidak valid."})
		return
	}

	now := time.Now()
	newSocialMedia := models.SocialMedia{
		Name:           payload.Name,
		SocialMediaURL: payload.SocialMediaURL,
		UserID:         currentUser.ID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	result := smc.DB.Create(&newSocialMedia)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":               newSocialMedia.ID,
		"name":             newSocialMedia.Name,
		"social_media_url": newSocialMedia.SocialMediaURL,
		"user_id":          newSocialMedia.UserID,
	})
}

// GetSocialMedias mengambil daftar media sosial yang dimiliki oleh pengguna saat ini
func (smc *SocialMediaController) GetSocialMedias(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var socialMedias []models.SocialMedia
	result := smc.DB.Preload("User").Where("user_id = ?", currentUser.ID).Find(&socialMedias)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	// Membuat slice untuk menyimpan respons yang sesuai dengan spesifikasi OpenAPI
	var responseData []gin.H
	for _, socialMedia := range socialMedias {
		responseData = append(responseData, gin.H{
			"id":               socialMedia.ID,
			"name":             socialMedia.Name,
			"social_media_url": socialMedia.SocialMediaURL,
			"user_id":          socialMedia.UserID,
			"user": gin.H{
				"id":       socialMedia.User.ID,
				"email":    socialMedia.User.Email,
				"username": socialMedia.User.Username,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// GetSocialMediaByID mengambil informasi media sosial berdasarkan ID
func (smc *SocialMediaController) GetSocialMediaByID(ctx *gin.Context) {
	socialMediaID := ctx.Param("socialMediaId")

	var socialMedia models.SocialMedia
	result := smc.DB.Preload("User").First(&socialMedia, "id = ?", socialMediaID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada social media dengan ID tersebut."})
		return
	}
	user := models.User{}
	smc.DB.First(&user, socialMedia.UserID)

	// Membuat objek JSON yang sesuai dengan spesifikasi OpenAPI
	responseData := gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"user": gin.H{
			"id":       socialMedia.User.ID,
			"email":    socialMedia.User.Email,
			"username": socialMedia.User.Username,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// UpdateSocialMedia digunakan untuk memperbarui media sosial yang ada
func (smc *SocialMediaController) UpdateSocialMedia(ctx *gin.Context) {
	socialMediaID := ctx.Param("socialMediaId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload models.UpdateSocialMediaRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validasi field name
	if payload.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Name harus diisi."})
		return
	}

	// Validasi field social_media_url
	if payload.SocialMediaURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "URL Social media harus diisi."})
		return
	}

	// Validasi URL profil
	if !utils.IsValidURL(payload.SocialMediaURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format profile image URL tidak valid."})
		return
	}

	var updatedSocialMedia models.SocialMedia
	result := smc.DB.First(&updatedSocialMedia, "id = ?", socialMediaID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada social media dengan ID tersebut."})
		return
	}

	// Periksa apakah pengguna yang sedang masuk adalah pemilik dari entri media sosial yang akan diperbarui
	if updatedSocialMedia.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Kamu tidak diizinkan untuk mengedit social media ini."})
		return
	}

	updatedSocialMedia.Name = payload.Name
	updatedSocialMedia.SocialMediaURL = payload.SocialMediaURL
	smc.DB.Save(&updatedSocialMedia)

	// Membuat objek JSON yang sesuai dengan spesifikasi OpenAPI
	responseData := gin.H{
		"id":               updatedSocialMedia.ID,
		"name":             updatedSocialMedia.Name,
		"social_media_url": updatedSocialMedia.SocialMediaURL,
		"user_id":          updatedSocialMedia.UserID,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// DeleteSocialMedia digunakan untuk menghapus entri media sosial
func (smc *SocialMediaController) DeleteSocialMedia(ctx *gin.Context) {
	socialMediaID := ctx.Param("socialMediaId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var socialMedia models.SocialMedia
	result := smc.DB.First(&socialMedia, "id = ?", socialMediaID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada social media dengan ID tersebut."})
		return
	}

	// Periksa apakah pengguna yang sedang masuk adalah pemilik dari entri media sosial yang akan dihapus
	if socialMedia.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Kamu tidak diizinkan untuk menghapus social media ini."})
		return
	}

	result = smc.DB.Delete(&socialMedia)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada social media dengan ID tersebut."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
