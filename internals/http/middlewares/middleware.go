package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/omniflare/go_starter/internals/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"error":  "unauthorized: no token provided",
			})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(cookie, "Bearer ")
		claims := jwt.MapClaims{}
		jwtSecret := config.NewEnvConfig().JWT_SECRET
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"error":  "unauthorized: invalid token",
			})
			c.Abort()
			return
		}
		userID := uint(claims["id"].(float64))
		userEmail := claims["email"].(string)
		c.Set("user_id", userID)
		c.Set("user_email", userEmail)

		c.Next()
	}
}
