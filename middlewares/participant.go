package middleware

import (
	"CfBE/initializers"
	"CfBE/models"
	"CfBE/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParticipantRoleRequired() gin.HandlerFunc {
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

		// Check if the user is a participa

		contestID := c.Param("contest_id")

		// Check if the contest exists
		var contest models.Contest

		// Check if the user is a participant
		if err := initializers.DB.Where("contest_id = ? AND Participants = ?", contestID, user.ID).First(&contest).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not a participant of this contest"})
			c.Abort()
			return
		}

		// If the user is an admin, continue with the request
		c.Next()
	}
}
