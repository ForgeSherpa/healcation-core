package profile

import (
	"healcationBackend/database"
	"healcationBackend/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.Select("id, username, email, created_at, updated_at").Where("id = ?", userID).First(&user).Error; err != nil {
		sendResponse(c, http.StatusNotFound, nil, "User not found")
		return
	}

	type UserResponse struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	data := UserResponse{
		ID:        user.ID,
		Name:      user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	sendResponse(c, http.StatusOK, data, "User profile retrieved successfully")
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		sendResponse(c, http.StatusNotFound, nil, "User not found")
		return
	}

	var data struct {
		Name     string `json:"name,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}

	if err := c.BindJSON(&data); err != nil {
		sendResponse(c, http.StatusBadRequest, nil, "Invalid request: "+err.Error())
		return
	}

	data.Email = strings.ToLower(data.Email)

	updates := map[string]interface{}{}

	if data.Name != "" {
		updates["username"] = data.Name
	}

	if data.Email != "" {
		updates["email"] = strings.ToLower(data.Email)
	}

	if data.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			sendResponse(c, http.StatusInternalServerError, nil, "Failed to encrypt new password")
			return
		}
		updates["password"] = string(hashed)
	}

	if len(updates) == 0 {
		sendResponse(c, http.StatusBadRequest, nil, "No fields to update")
		return
	}

	if err := database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(updates).Error; err != nil {
		sendResponse(c, http.StatusInternalServerError, nil, "Failed to update profile: "+err.Error())
		return
	}

	sendResponse(c, http.StatusOK, nil, "Profile updated successfully")
}
