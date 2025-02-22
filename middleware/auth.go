package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strconv"
)

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token provided"})
			c.Abort()
			return
		}
		fmt.Println("Access Token from Cookie: ", tokenString) // Debugging line

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Handle both string and float64 for "sub"
			userID := claims["sub"]
			switch v := userID.(type) {
			case string:
				c.Set("userID", v)
			case float64:
				c.Set("userID", strconv.FormatFloat(v, 'f', -1, 64))
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID"})
				c.Abort()
				return
			}
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
		}
	}
}
