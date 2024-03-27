package routes

import (
	"mygram-final-project/controllers"
	"mygram-final-project/middleware"

	"github.com/gin-gonic/gin"
)

// PhotoRouteController mengelola rute yang terkait dengan foto.
type PhotoRouteController struct {
	photoController controllers.PhotoController // Kontroler untuk foto
}

// NewRoutePhotoController membuat instance baru dari PhotoRouteController.
func NewRoutePhotoController(photoController controllers.PhotoController) PhotoRouteController {
	return PhotoRouteController{photoController}
}

// PhotoRoute menentukan rute yang terkait dengan foto.
func (pc *PhotoRouteController) PhotoRoute(rg *gin.RouterGroup) {
	router := rg.Group("photos")
	router.Use(middleware.UserExtractor())

	router.POST("", pc.photoController.CreatePhoto)            // Rute untuk membuat foto baru
	router.GET("", pc.photoController.FindPhotos)              // Rute untuk menemukan semua foto
	router.PUT("/:photoId", pc.photoController.UpdatePhoto)    // Rute untuk memperbarui foto berdasarkan ID
	router.GET("/:photoId", pc.photoController.FindPhotoByID)  // Rute untuk menemukan foto berdasarkan ID
	router.DELETE("/:photoId", pc.photoController.DeletePhoto) // Rute untuk menghapus foto berdasarkan ID
}
