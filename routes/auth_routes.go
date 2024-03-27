package routes

import (
	"mygram-final-project/controllers"

	"github.com/gin-gonic/gin"
)

// AuthRouteController mengelola rute yang terkait dengan otentikasi pengguna.
type AuthRouteController struct {
	authController controllers.AuthController // Kontroler untuk otentikasi pengguna
	userController controllers.UserController // Kontroler untuk pengguna
}

// NewAuthRouteController membuat instance baru dari AuthRouteController.
func NewAuthRouteController(authController controllers.AuthController, userController controllers.UserController) AuthRouteController {
	return AuthRouteController{authController, userController} // Inisialisasi kontroler pengguna
}

// AuthRoute menentukan rute yang terkait dengan otentikasi pengguna.
func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
}
