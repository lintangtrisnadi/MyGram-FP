package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mygram-final-project/initializers"
	"mygram-final-project/models"
	"mygram-final-project/utils"
)

// AuthController adalah kontroler untuk otentikasi pengguna
type AuthController struct {
	DB *gorm.DB
}

// NewAuthController digunakan untuk membuat instance baru dari AuthController
func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// SignUpUser digunakan untuk menangani permintaan pendaftaran pengguna baru
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload models.SignUpInput

	// Bind JSON payload ke struct SignUpInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Validasi email
	if !utils.IsValidEmail(payload.Email) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format email tidak valid."})
		return
	}

	// Validasi email unik
	var existingUser models.User
	if ac.DB.Where("email = ?", strings.ToLower(payload.Email)).First(&existingUser).RowsAffected != 0 {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User dengan email tersebut sudah ada."})
		return
	}

	// Validasi username
	if payload.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Username harus diisi."})
		return
	}
	if ac.DB.Where("username = ?", payload.Username).First(&existingUser).RowsAffected != 0 {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Username tersebut sudah ada."})
		return
	}

	// Validasi password
	if len(payload.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Password setidaknya berisi 6 karakter."})
		return
	}

	// Validasi age
	if payload.Age < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Umur minimal 8 tahun."})
		return
	}

	// Validasi URL profil
	if payload.ProfileImageURL != "" && !utils.IsValidURL(payload.ProfileImageURL) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Format URL profile image tidak valid."})
		return
	}

	// Hash password sebelum menyimpan ke database
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Username:        payload.Username,
		Email:           strings.ToLower(payload.Email),
		Password:        hashedPassword,
		Age:             payload.Age,
		ProfileImageURL: payload.ProfileImageURL,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Simpan pengguna baru ke database
	result := ac.DB.Create(&newUser)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	// Respon dengan data pengguna yang baru dibuat
	userResponse := gin.H{
		"id":                newUser.ID,
		"email":             newUser.Email,
		"username":          newUser.Username,
		"age":               newUser.Age,
		"profile_image_url": newUser.ProfileImageURL,
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": userResponse})
}

// SignInUser digunakan untuk menangani permintaan masuk pengguna
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload models.SignInInput

	// Bind JSON payload ke struct SignInInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Cari pengguna berdasarkan alamat email
	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Email atau password tidak valid."})
		return
	}

	// Verifikasi password
	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Email atau password tidak valid."})
		return
	}

	// Load konfigurasi dari file .env
	config, _ := initializers.LoadConfig(".")

	// Generate token akses
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Generate token refresh
	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Set cookie token pada response
	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)

	// Respon dengan token akses
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}
