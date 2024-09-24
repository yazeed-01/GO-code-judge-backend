package middleware

import (
	"CfBE/initializers"
	"CfBE/models"
	"CfBE/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRoleRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the access token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Parse and validate the token
		claims, err := utils.ParseToken(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Find the user in the database
		var user models.User
		if err := initializers.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Check if the user has an admin role
		if user.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		// If the user is an admin, continue with the request
		c.Next()
	}
}
