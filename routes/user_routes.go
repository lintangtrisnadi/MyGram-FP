package routes

import (
	"mygram-final-project/controllers"
	"mygram-final-project/middleware"

	"github.com/gin-gonic/gin"
)

// UserRouteController mengelola rute terkait pengguna.
type UserRouteController struct {
	userController controllers.UserController // Kontroler untuk pengguna
}

// NewRouteUserController membuat instance baru dari UserRouteController.
func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

// UserRoute menentukan rute terkait pengguna.
func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")

	// Mengatur rute untuk memperbarui informasi pengguna
	router.PUT("", middleware.UserExtractor(), uc.userController.UpdateMe)
	// Mengatur rute untuk menghapus pengguna
	router.DELETE("", middleware.UserExtractor(), uc.userController.DeleteMe)
}
