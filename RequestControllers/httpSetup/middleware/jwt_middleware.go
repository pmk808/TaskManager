package middleware

import (
	"net/http"
	"strings"
	"taskmanager/RequestControllers/httpSetup/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization format",
			})
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		c.Set("client_name", claims.ClientName)
		c.Set("client_id", claims.ClientID)

		c.Next()
	}
}
