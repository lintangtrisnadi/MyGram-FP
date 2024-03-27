package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mygram-final-project/models"
	"mygram-final-project/utils"
)

// PhotoController adalah kontroler untuk pengelolaan foto
type PhotoController struct {
	DB *gorm.DB
}

// NewPhotoController digunakan untuk membuat instance baru dari PhotoController
func NewPhotoController(DB *gorm.DB) PhotoController {
	return PhotoController{DB}
}

// CreatePhoto digunakan untuk menangani permintaan pembuatan foto baru
func (pc *PhotoController) CreatePhoto(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreatePhotoRequest

	// Bind JSON payload ke struct CreatePhotoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validasi field title
	if payload.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Title harus diisi."})
		return
	}

	// Validasi field photo_url
	if payload.PhotoURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Photo URL harus diisi."})
		return
	}

	// Validasi URL foto
	if payload.PhotoURL != "" && !utils.IsValidURL(payload.PhotoURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format photo URL tidak valid"})
		return
	}

	now := time.Now()
	newPhoto := models.Photo{
		Title:     payload.Title,
		Caption:   payload.Caption,
		PhotoURL:  payload.PhotoURL,
		UserID:    currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := pc.DB.Create(&newPhoto)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"message": "Photo dengan title tersebut sudah ada."})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":        newPhoto.ID,
		"caption":   newPhoto.Caption,
		"title":     newPhoto.Title,
		"photo_url": newPhoto.PhotoURL,
		"user_id":   newPhoto.UserID,
	})
}

// UpdatePhoto digunakan untuk menangani permintaan pembaruan foto
func (pc *PhotoController) UpdatePhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdatePhoto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Validasi field title
	if payload.Title == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Title harus diisi."})
		return
	}

	// Validasi field photo_url
	if payload.PhotoURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Photo URL harus diisi."})
		return
	}

	// Validasi URL foto
	if !utils.IsValidURL(payload.PhotoURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format photo URL tidak valid."})
		return
	}

	var updatedPhoto models.Photo
	result := pc.DB.First(&updatedPhoto, "id = ?", photoID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Tidak ada photo dengan ID tersebut."})
		return
	}

	// Periksa apakah pengguna yang sedang masuk adalah pemilik foto yang akan diperbarui
	if updatedPhoto.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Kamu tidak diizinkan untuk mengedit photo ini."})
		return
	}

	now := time.Now()
	updatedPhoto.Title = payload.Title
	updatedPhoto.Caption = payload.Caption
	updatedPhoto.PhotoURL = payload.PhotoURL
	updatedPhoto.UpdatedAt = now

	pc.DB.Save(&updatedPhoto)

	// Mengonversi data yang diperbarui menjadi respons sesuai dengan spesifikasi OpenAPI
	responseData := gin.H{
		"id":        updatedPhoto.ID,
		"caption":   updatedPhoto.Caption,
		"title":     updatedPhoto.Title,
		"photo_url": updatedPhoto.PhotoURL,
		"user_id":   updatedPhoto.UserID,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// FindPhotoByID digunakan untuk menemukan foto berdasarkan ID
func (pc *PhotoController) FindPhotoByID(ctx *gin.Context) {
	photoID := ctx.Param("photoId")

	var photo models.Photo
	result := pc.DB.First(&photo, "id = ?", photoID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Tidak ada photo dengan ID tersebut."})
		return
	}

	// Ambil informasi pengguna yang terkait dengan foto dari basis data
	user := models.User{}
	pc.DB.First(&user, photo.UserID)

	// Membuat objek JSON yang sesuai dengan spesifikasi OpenAPI
	responseData := gin.H{
		"id":        photo.ID,
		"caption":   photo.Caption,
		"title":     photo.Title,
		"photo_url": photo.PhotoURL,
		"user_id":   photo.UserID,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// FindPhotos digunakan untuk menemukan daftar foto dengan opsi paging
func (pc *PhotoController) FindPhotos(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var photos []models.Photo
	results := pc.DB.Limit(intLimit).Offset(offset).Find(&photos)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	// Membuat slice untuk menyimpan hasil response yang sesuai dengan spesifikasi OpenAPI
	var responseData []gin.H
	for _, photo := range photos {
		user := models.User{}            // Deklarasikan variabel untuk menyimpan informasi pengguna
		pc.DB.First(&user, photo.UserID) // Ambil informasi pengguna dari basis data berdasarkan ID yang terkait dengan foto

		responseData = append(responseData, gin.H{
			"id":        photo.ID,
			"caption":   photo.Caption,
			"title":     photo.Title,
			"photo_url": photo.PhotoURL,
			"user_id":   photo.UserID,
			"user": gin.H{
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// DeletePhoto digunakan untuk menghapus foto
func (pc *PhotoController) DeletePhoto(ctx *gin.Context) {
	photoID := ctx.Param("photoId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var photo models.Photo
	result := pc.DB.First(&photo, "id = ?", photoID)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Tidak ada photo dengan ID tersebut."})
		return
	}

	// Periksa apakah pengguna yang masuk adalah pemilik foto
	if photo.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "Kamu tidak diizinkan untuk menghapus photo ini."})
		return
	}

	// Hapus foto dari basis data
	result = pc.DB.Delete(&photo)

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Tidak ada photo dengan ID tersebut."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
