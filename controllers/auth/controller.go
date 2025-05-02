package auth

import (
	"healcationBackend/database"
	"healcationBackend/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func sendResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, Response{
		Status:  status,
		Data:    data,
		Message: message,
	})
}

func Register(c *gin.Context) {
	var data struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Format JSON salah")
		return
	}

	data.Username = strings.ToLower(data.Username)
	data.Email = strings.ToLower(data.Email)

	var existing models.User
	if err := database.DB.Where("email = ?", data.Email).First(&existing).Error; err == nil {
		sendResponse(c, http.StatusConflict, nil, "Email sudah terdaftar")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal mengenkripsi password")
		return
	}

	user := models.User{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hash),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal membuat pengguna baru")
		return
	}

	sendResponse(c, http.StatusOK, nil, "Registrasi berhasil")
}

func Login(c *gin.Context) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Format JSON salah")
		return
	}

	data.Email = strings.ToLower(data.Email)

	var user models.User
	if err := database.DB.First(&user, "email = ?", data.Email).Error; err != nil {
		sendResponse(c, http.StatusUnauthorized, nil, "Email atau password salah")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		sendResponse(c, http.StatusUnauthorized, nil, "Email atau password salah")
		return
	}
	env := os.Getenv("APP_ENV") // "staging" or "production"
	var accessExp, refreshExp time.Time
	if env == "staging" {
		accessExp = time.Now().Add(2 * time.Hour)
		refreshExp = time.Now().Add(2 * time.Hour)
	} else {
		accessExp = time.Now().Add(7 * 24 * time.Hour)
		refreshExp = time.Now().Add(30 * 24 * time.Hour)
	}

	userIDStr := strconv.FormatUint(uint64(user.ID), 10)
	secret := []byte(os.Getenv("SECRET"))

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userIDStr,
		"exp": accessExp.Unix(),
	}).SignedString(secret)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal membuat Access Token")
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userIDStr,
		"exp": refreshExp.Unix(),
	}).SignedString(secret)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal membuat Refresh Token")
		return
	}

	sendResponse(c, http.StatusOK, gin.H{
		"access_token":             accessToken,
		"access_token_expired_at":  accessExp.Format(time.RFC3339),
		"refresh_token":            refreshToken,
		"refresh_token_expired_at": refreshExp.Format(time.RFC3339),
	}, "Login berhasil")
}

func RefreshToken(c *gin.Context) {
	var data TokenRequest

	if err := c.ShouldBindJSON(&data); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Format JSON salah")
		return
	}

	secret := []byte(os.Getenv("SECRET"))
	refreshToken, err := jwt.Parse(data.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})

	if err != nil || !refreshToken.Valid {
		sendResponse(c, http.StatusUnauthorized, nil, "Refresh Token tidak valid atau kadaluarsa")
		return
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil {
		sendResponse(c, http.StatusUnauthorized, nil, "Refresh Token tidak memiliki klaim yang valid")
		return
	}

	userID := claims["sub"].(string)
	expAccess := time.Now().Add(time.Hour * 24 * 7)
	expRefresh := time.Now().Add(time.Hour * 24 * 30)

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": expAccess.Unix(),
	}).SignedString(secret)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal membuat Access Token baru")
		return
	}

	newRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": expRefresh.Unix(),
	}).SignedString(secret)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Gagal membuat Refresh Token baru")
		return
	}

	sendResponse(c, http.StatusOK, gin.H{
		"access_token":             accessToken,
		"access_token_expired_at":  expAccess,
		"refresh_token":            newRefreshToken,
		"refresh_token_expired_at": expRefresh,
	}, "Token berhasil diperbarui")
}

func Validate(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	data := gin.H{"userID": userID}
	sendResponse(c, http.StatusOK, data, "User profile retrieved successfully")
}
