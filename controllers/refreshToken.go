package controllers

import (
	"CfBE/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefreshToken handles the refreshing of access tokens
func RefreshToken(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	// Use the token directly since you're not using "Bearer"
	token := authHeader
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Get the user ID from the claims
	userID := claims.UserID // Ensure claims.UserID is of the correct type (uint)

	// Generate new tokens
	newAccessToken, _, err := utils.GenerateTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	// Return the new access token
	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
