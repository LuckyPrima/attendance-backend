package controllers

import (
	"net/http"
	"time"

	"github.com/LuckyPrima/attendance-backend/config"
	"github.com/LuckyPrima/attendance-backend/models"
	"github.com/gin-gonic/gin"
)

// helper parse time string (support HH:mm dan HH:mm:ss)
func parseTimeFlexible(val string) error {
	layoutShort := "15:04"
	layoutFull := "15:04:05"

	if _, err := time.Parse(layoutShort, val); err == nil {
		return nil
	}
	if _, err := time.Parse(layoutFull, val); err == nil {
		return nil
	}
	return &time.ParseError{}
}

// ================== CREATE ==================
func CreateDepartment(c *gin.Context) {
	var input struct {
		DepartementName string `json:"departement_name"`
		MaxClockIn      string `json:"max_clock_in"`
		MaxClockOut     string `json:"max_clock_out"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation for time format
	if err := parseTimeFlexible(input.MaxClockIn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for max_clock_in, use HH:mm or HH:mm:ss"})
		return
	}
	if err := parseTimeFlexible(input.MaxClockOut); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for max_clock_out, use HH:mm or HH:mm:ss"})
		return
	}

	department := models.Department{
		DepartementName: input.DepartementName,
		MaxClockIn:      input.MaxClockIn,
		MaxClockOut:     input.MaxClockOut,
	}

	if err := config.DB.Create(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	c.JSON(http.StatusCreated, department)
}

// ================== READ ==================
func GetDepartments(c *gin.Context) {
	var departments []models.Department

	if err := config.DB.Find(&departments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments"})
		return
	}

	c.JSON(http.StatusOK, departments)
}

// ================== UPDATE ==================
func UpdateDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.First(&department, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	var input struct {
		DepartementName *string `json:"departement_name"`
		MaxClockIn      *string `json:"max_clock_in"`
		MaxClockOut     *string `json:"max_clock_out"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.DepartementName != nil {
		department.DepartementName = *input.DepartementName
	}
	if input.MaxClockIn != nil {
		if err := parseTimeFlexible(*input.MaxClockIn); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for max_clock_in, use HH:mm or HH:mm:ss"})
			return
		}
		department.MaxClockIn = *input.MaxClockIn
	}
	if input.MaxClockOut != nil {
		if err := parseTimeFlexible(*input.MaxClockOut); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format for max_clock_out, use HH:mm or HH:mm:ss"})
			return
		}
		department.MaxClockOut = *input.MaxClockOut
	}

	if err := config.DB.Save(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department"})
		return
	}

	c.JSON(http.StatusOK, department)
}

// ================== DELETE ==================
func DeleteDepartment(c *gin.Context) {
	id := c.Param("id")
	var department models.Department

	if err := config.DB.First(&department, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	if err := config.DB.Delete(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete department"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully"})
}
