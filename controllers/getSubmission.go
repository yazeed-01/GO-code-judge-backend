package controllers

import (
	"CfBE/initializers"
	"CfBE/models"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetSubmission(c *gin.Context) {
	// Take ID from URL
	resultIDStr := c.Param("id4")
	var resultID uint64
	var err error

	// Convert ID from string to uint
	if resultID, err = strconv.ParseUint(resultIDStr, 10, 32); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid result ID"})
		return
	}

	var result models.Result
	// Fetch the result using ID
	if err := initializers.DB.First(&result, resultID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	// Fetch user details based on UserID
	var user models.User // Assume you have a User model defined
	if err := initializers.DB.First(&user, result.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Decode the results
	OutputResultDecode, err := base64.StdEncoding.DecodeString(result.Output)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding output"})
		return
	}

	// Convert to string
	decodedOutputString := string(OutputResultDecode)

	// Return only Username and FullName
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"username":  user.Username, // Assuming you have a Username field in your User model
			"full_name": user.FullName, // Assuming you have a FullName field in your User model
		},
		"submission": gin.H{
			"id":                 result.ID, // Use ID here
			"status_id":          result.StatusID,
			"status_description": result.StatusDescription,
			"memory_used":        result.MemoryUsed,
			"created_at":         result.CreatedAt,
			"finished_at":        result.FinishedAt,
			"language_id":        result.LanguageID,
			"source_code":        result.SourceCode,
			"input":              result.Input,
			"output":             decodedOutputString,
			"error_output":       result.ErrorOutput,
			"result_message":     result.ResultMessage,
		},
	})
}
