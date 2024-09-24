package controllers

import (
	"CfBE/initializers"
	"CfBE/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateProblem(c *gin.Context) {
	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contestIDStr := c.Param("id2")                            // Get contest ID from URL
	contestID, err := strconv.ParseUint(contestIDStr, 10, 32) // Convert string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}
	problem.ContestID = uint(contestID) // Assign contest ID to the problem

	// Create the problem
	if err := initializers.DB.Create(&problem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"problem": problem})
}

func DeleteProblem(c *gin.Context) {
	contestIDStr := c.Param("id2") // Get contest ID from URL
	problemID := c.Param("id3")    // Get problem ID from URL

	contestID, err := strconv.ParseUint(contestIDStr, 10, 32) // Convert string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}

	var problem models.Problem
	// Find the problem by ID
	if err := initializers.DB.First(&problem, problemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	// Check if the problem belongs to the specified contest
	if problem.ContestID != uint(contestID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Problem does not belong to the specified contest"})
		return
	}

	// Delete the problem
	if err := initializers.DB.Delete(&problem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting problem"})
		return
	}

	c.Status(http.StatusOK)
}

func Problems(c *gin.Context) {
	contestIDStr := c.Param("id2") // Get contest ID from URL

	contestID, err := strconv.ParseUint(contestIDStr, 10, 32) // Convert string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}

	var problems []models.Problem
	// Fetch all problems related to the given contest
	if err := initializers.DB.Where("contest_id = ?", uint(contestID)).Find(&problems).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

func Problem(c *gin.Context) {
	contestIDStr := c.Param("id2") // Get contest ID from URL
	problemID := c.Param("id3")    // Get problem ID from URL

	contestID, err := strconv.ParseUint(contestIDStr, 10, 32) // Convert string to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}

	var problem models.Problem
	// Find the problem by ID
	if err := initializers.DB.First(&problem, problemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Problem not found"})
		return
	}

	// Check if the problem belongs to the specified contest
	if problem.ContestID != uint(contestID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Problem does not belong to the specified contest"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"problem": problem})
}
