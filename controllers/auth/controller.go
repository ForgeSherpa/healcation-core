package auth

import (
	"healcationBackend/database"
	"healcationBackend/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var data struct {
		Username string
		Email    string
		Password string
	}

	if c.Bind(&data) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read data"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: data.Username,
		Email:    data.Email,
		Password: string(hash),
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func Login(c *gin.Context) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read data"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", data.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email or Password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Email or Password"})
		return
	}

	expAccess := time.Now().Add(time.Hour * 24 * 7)
	expRefresh := time.Now().Add(time.Hour * 24 * 30)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": expAccess.Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Access Token"})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": expRefresh.Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Refresh Token"})
		return
	}

	user.RefreshToken = refreshTokenString
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"access_token":             accessTokenString,
		"access_token_expired_at":  expAccess,
		"refresh_token":            refreshTokenString,
		"refresh_token_expired_at": expRefresh,
	})
}

func Validate(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userID": userID})
}

func Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if exists {
		database.DB.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", "")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func RefreshToken(c *gin.Context) {
	var data struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read data"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "refresh_token = ?", data.RefreshToken).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired Refresh Token"})
		return
	}

	expAccess := time.Now().Add(time.Hour * 24 * 7)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": expAccess.Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Access Token"})
		return
	}

	database.DB.Model(&user).Update("refresh_token", "")

	c.JSON(http.StatusOK, gin.H{
		"access_token":            accessTokenString,
		"access_token_expired_at": expAccess,
	})
}
