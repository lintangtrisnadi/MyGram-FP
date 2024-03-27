package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// IsValidEmail memeriksa apakah alamat email yang diberikan berada dalam format yang valid.
func IsValidEmail(email string) bool {
	// Pola regex untuk validasi email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// HashPassword menghasilkan hash dari password yang diberikan menggunakan bcrypt.
func HashPassword(password string) (string, error) {
	// Generate hash dari password dengan menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

// VerifyPassword memverifikasi apakah password kandidat cocok dengan hash password yang diberikan.
func VerifyPassword(hashedPassword string, candidatePassword string) error {
	// Membandingkan password kandidat dengan hash password menggunakan bcrypt
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}

// CreateToken membuat token JWT dengan menggunakan private key yang diberikan.
func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	// Mendecode private key dari base64
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	// Parse private key RSA dari PEM
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	// Menyiapkan klaim token JWT
	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	// Membuat token JWT dengan klaim yang disiapkan sebelumnya
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

// ValidateToken memvalidasi token JWT menggunakan public key yang diberikan.
func ValidateToken(token string, publicKey string) (interface{}, error) {
	// Mendecode public key dari base64
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	// Parse public key RSA dari PEM
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	// Memvalidasi token JWT dengan menggunakan public key yang sudah diparse
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	// Mendapatkan klaim dari token yang sudah divalidasi
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}

// IsValidURL memeriksa apakah string URL yang diberikan valid atau tidak.
func IsValidURL(urlString string) bool {
	// ParseRequestURI digunakan untuk memeriksa apakah string URL dapat di-parse ke URL yang valid
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		// Jika terjadi kesalahan saat parsing, URL dianggap tidak valid
		return false
	}
	// Jika tidak ada kesalahan, URL dianggap valid
	return true
}
