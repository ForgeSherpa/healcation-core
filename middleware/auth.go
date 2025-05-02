package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendResponse(c, http.StatusUnauthorized, nil, "Invalid token format")
			c.Abort()
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil || !token.Valid {
			sendResponse(c, http.StatusUnauthorized, nil, "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			sendResponse(c, http.StatusUnauthorized, nil, "Unauthorized")
			c.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			sendResponse(c, http.StatusUnauthorized, nil, "Invalid or missing sub claim")
			c.Abort()
			return
		}

		id64, err := strconv.ParseUint(sub, 10, 64)
		if err != nil {
			sendResponse(c, http.StatusUnauthorized, nil, "Invalid sub format")
			c.Abort()
			return
		}
		c.Set("userID", uint(id64))

		c.Next()
	}
}
