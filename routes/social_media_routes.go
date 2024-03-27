package routes

import (
	"mygram-final-project/controllers"
	"mygram-final-project/middleware"

	"github.com/gin-gonic/gin"
)

// SocialMediaRouteController mengelola rute yang terkait dengan media sosial.
type SocialMediaRouteController struct {
	socialMediaController controllers.SocialMediaController // Kontroler untuk media sosial
}

// NewRouteSocialMediaController membuat instance baru dari SocialMediaRouteController.
func NewRouteSocialMediaController(socialMediaController controllers.SocialMediaController) SocialMediaRouteController {
	return SocialMediaRouteController{socialMediaController}
}

// SocialMediaRoute menentukan rute yang terkait dengan media sosial.
func (smc *SocialMediaRouteController) SocialMediaRoute(rg *gin.RouterGroup) {
	router := rg.Group("socialmedias")
	router.Use(middleware.UserExtractor())

	router.POST("", smc.socialMediaController.CreateSocialMedia)                  // Rute untuk membuat media sosial baru
	router.GET("", smc.socialMediaController.GetSocialMedias)                     // Rute untuk mendapatkan semua media sosial
	router.GET("/:socialMediaId", smc.socialMediaController.GetSocialMediaByID)   // Rute untuk mendapatkan media sosial berdasarkan ID
	router.PUT("/:socialMediaId", smc.socialMediaController.UpdateSocialMedia)    // Rute untuk memperbarui media sosial berdasarkan ID
	router.DELETE("/:socialMediaId", smc.socialMediaController.DeleteSocialMedia) // Rute untuk menghapus media sosial berdasarkan ID
}
