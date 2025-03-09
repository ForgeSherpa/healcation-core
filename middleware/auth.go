package middleware

import (
	"fmt"
	"net/http"
	"os"

	// "strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func sendResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"message": message,
		"data":    data,
	})
}

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sendResponse(c, http.StatusUnauthorized, nil, "No access token provided")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			sendResponse(c, http.StatusUnauthorized, nil, "Invalid token format")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			sendResponse(c, http.StatusUnauthorized, nil, "Invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized")
			c.Abort()
			return
		}

		userID, exists := claims["sub"]
		if !exists {
			sendResponse(c, http.StatusUnauthorized, nil, "User ID not found in token")
			c.Abort()
			return
		}

		// Menyimpan userID dalam context dengan tipe uint
		// var userIDUint uint
		// switch v := userID.(type) {
		// case float64: // Jika userID berupa float64 (misalnya dari JSON parsing)
		// 	userIDUint = uint(v)
		// case string: // Jika userID berupa string, konversi ke uint
		// 	idInt, err := strconv.ParseUint(v, 10, 64)
		// 	if err != nil {
		// 		sendResponse(c, http.StatusUnauthorized, nil, "Invalid user ID format")
		// 		c.Abort()
		// 		return
		// 	}
		// 	userIDUint = uint(idInt)
		// default:
		// 	sendResponse(c, http.StatusUnauthorized, nil, "Invalid user ID type")
		// 	c.Abort()
		// 	return
		// }

		// c.Set("userID", userIDUint)
		// c.Next()

		// Menyimpan userID dalam context dengan tipe string
		switch v := userID.(type) {
		case string:
			c.Set("userID", v)
		case float64:
			c.Set("userID", fmt.Sprintf("%.0f", v))
		default:
			sendResponse(c, http.StatusUnauthorized, nil, "Invalid user ID format")
			c.Abort()
			return
		}
	}
}
