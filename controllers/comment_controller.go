package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mygram-final-project/models"
)

// CommentController adalah kontroler untuk operasi yang berhubungan dengan komentar
type CommentController struct {
	DB *gorm.DB
}

// NewCommentController digunakan untuk membuat instance baru dari CommentController
func NewCommentController(DB *gorm.DB) CommentController {
	return CommentController{DB}
}

// CreateComment digunakan untuk menangani permintaan pembuatan komentar baru
func (cc *CommentController) CreateComment(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateCommentRequest

	// Bind JSON payload ke struct CreateCommentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validasi field message
	if payload.Message == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Message harus diisi."})
		return
	}

	now := time.Now()
	newComment := models.Comment{
		Message:   payload.Message,
		PhotoID:   payload.PhotoID,
		UserID:    currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Simpan komentar baru ke database
	result := cc.DB.Create(&newComment)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":       newComment.ID,
		"message":  newComment.Message,
		"photo_id": newComment.PhotoID,
		"user_id":  newComment.UserID,
	})
}

// GetComments digunakan untuk menangani permintaan untuk mendapatkan semua komentar
func (cc *CommentController) GetComment(ctx *gin.Context) {
	var comments []models.Comment
	result := cc.DB.Preload("User").Preload("Photo").Find(&comments)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Unexpected error"})
		return
	}

	// Membuat slice untuk menyimpan respons yang sesuai dengan spesifikasi OpenAPI
	var responseData []gin.H
	for _, comment := range comments {
		responseData = append(responseData, gin.H{
			"id":       comment.ID,
			"message":  comment.Message,
			"photo_id": comment.PhotoID,
			"user_id":  comment.UserID,
			"user": gin.H{
				"id":       comment.User.ID,
				"email":    comment.User.Email,
				"username": comment.User.Username,
			},
			"photo": gin.H{
				"id":        comment.Photo.ID,
				"caption":   comment.Photo.Caption,
				"title":     comment.Photo.Title,
				"photo_url": comment.Photo.PhotoURL,
				"user_id":   comment.Photo.UserID,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// GetCommentByID digunakan untuk menangani permintaan untuk mendapatkan komentar berdasarkan ID
func (cc *CommentController) GetCommentByID(ctx *gin.Context) {
	commentID := ctx.Param("commentId")

	var comment models.Comment
	result := cc.DB.Preload("User").Preload("Photo").First(&comment, "id = ?", commentID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Tidak ada komentar dengan ID tersebut."})
		return
	}

	// Membuat objek JSON yang sesuai dengan spesifikasi OpenAPI
	responseData := gin.H{
		"id":       comment.ID,
		"message":  comment.Message,
		"photo_id": comment.PhotoID,
		"user_id":  comment.UserID,
		"user": gin.H{
			"id":       comment.User.ID,
			"email":    comment.User.Email,
			"username": comment.User.Username,
		},
		"photo": gin.H{
			"id":        comment.Photo.ID,
			"caption":   comment.Photo.Caption,
			"title":     comment.Photo.Title,
			"photo_url": comment.Photo.PhotoURL,
			"user_id":   comment.Photo.UserID,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// UpdateComment digunakan untuk menangani permintaan untuk memperbarui komentar
func (cc *CommentController) UpdateComment(ctx *gin.Context) {
	commentID := ctx.Param("commentId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload models.UpdateCommentRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Validasi field message
	if payload.Message == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Message harus diisi."})
		return
	}

	var updatedComment models.Comment
	result := cc.DB.First(&updatedComment, "id = ?", commentID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada komentar dengan ID tersebut."})
		return
	}

	// Periksa apakah pengguna yang sedang masuk adalah pemilik komentar yang akan diperbarui
	if updatedComment.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Kamu tidak diizinkan untuk mengedit komentar ini."})
		return
	}

	updatedComment.Message = payload.Message
	cc.DB.Save(&updatedComment)

	// Membuat objek JSON yang sesuai dengan spesifikasi OpenAPI
	responseData := gin.H{
		"id":       updatedComment.ID,
		"message":  updatedComment.Message,
		"photo_id": updatedComment.PhotoID,
		"user_id":  updatedComment.UserID,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": responseData})
}

// DeleteComment digunakan untuk menangani permintaan untuk menghapus komentar
func (cc *CommentController) DeleteComment(ctx *gin.Context) {
	commentID := ctx.Param("commentId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var comment models.Comment
	result := cc.DB.First(&comment, "id = ?", commentID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada komentar dengan ID tersebut."})
		return
	}

	// Periksa apakah pengguna yang sedang masuk adalah pemilik komentar yang akan dihapus
	if comment.UserID != currentUser.ID {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Kamu tidak diizinkan untuk menghapus komentar ini."})
		return
	}

	result = cc.DB.Delete(&comment)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada komentar dengan ID tersebut."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
