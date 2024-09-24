package controllers

import (
	"CfBE/initializers"
	"CfBE/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Users(c *gin.Context) {
	var users []models.User
	if err := initializers.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract just the usernames
	var usernames []string
	for _, user := range users {
		usernames = append(usernames, user.Username)
	}

	c.JSON(http.StatusOK, gin.H{
		"usernames": usernames,
	})
}

func User(c *gin.Context) {
	var id = c.Param("id1") // Get the ID from the URL parameter
	var user models.User

	// Find the user using the ID rather than UserID
	if err := initializers.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return only Username and FullName
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"username": user.Username,
			"fullname": user.FullName,
		},
	})
}
