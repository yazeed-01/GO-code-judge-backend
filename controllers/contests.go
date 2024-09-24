package controllers

import (
	"CfBE/initializers"
	"CfBE/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateContest(c *gin.Context) {
	var contest models.Contest
	if err := c.ShouldBindJSON(&contest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := initializers.DB.Create(&contest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contest": contest,
	})
}

func DeleteContest(c *gin.Context) {
	id := c.Param("id2")
	initializers.DB.Delete(&models.Contest{}, id)
	c.Status(http.StatusOK)
}

func Contests(c *gin.Context) {
	var contests []models.Contest
	if err := initializers.DB.Find(&contests).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"contests": contests,
	})
}

func Contest(c *gin.Context) {
	var id = c.Param("id2")
	var contest models.Contest
	// find using UserID
	if err := initializers.DB.Where("contest_id = ?", id).First(&contest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contest": contest,
	})
}
