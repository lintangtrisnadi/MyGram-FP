package routes

import (
	"mygram-final-project/controllers"
	"mygram-final-project/middleware"

	"github.com/gin-gonic/gin"
)

// CommentRouteController mengelola rute yang terkait dengan komentar.
type CommentRouteController struct {
	commentController controllers.CommentController // Kontroler untuk komentar
}

// NewRouteCommentController membuat instance baru dari CommentRouteController.
func NewRouteCommentController(commentController controllers.CommentController) CommentRouteController {
	return CommentRouteController{commentController}
}

// CommentRoute menentukan rute yang terkait dengan komentar.
func (cc *CommentRouteController) CommentRoute(rg *gin.RouterGroup) {
	router := rg.Group("comments")
	router.Use(middleware.UserExtractor())

	router.POST("", cc.commentController.CreateComment)              // Rute untuk membuat komentar
	router.GET("", cc.commentController.GetComment)                  // Rute untuk mendapatkan semua komentar
	router.PUT("/:commentId", cc.commentController.UpdateComment)    // Rute untuk memperbarui komentar berdasarkan ID
	router.GET("/:commentId", cc.commentController.GetCommentByID)   // Rute untuk mendapatkan komentar berdasarkan ID
	router.DELETE("/:commentId", cc.commentController.DeleteComment) // Rute untuk menghapus komentar berdasarkan ID
}
