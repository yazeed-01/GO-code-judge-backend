package middleware

import (
	"CfBE/initializers"
	"CfBE/models"
	"CfBE/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SubmitAuth() gin.HandlerFunc {
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

		// Get the problem ID and contest ID from the URL parameters
		problemID := c.Param("id")
		contestID := c.Param("contestID")

		// Check if the contest exists
		var contest models.Contest
		if err := initializers.DB.First(&contest, contestID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contest not found"})
			c.Abort()
			return
		}

		// Check if the problem exists and is associated with the contest
		var problem models.Problem
		if err := initializers.DB.First(&problem, problemID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
			c.Abort()
			return
		}

		// Verify that the problem is linked to the contest
		if problem.ContestID != contest.ContestID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Problem does not belong to the specified contest"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}
