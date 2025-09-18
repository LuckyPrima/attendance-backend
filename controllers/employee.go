package controllers

import (
	"net/http"

	"github.com/LuckyPrima/attendance-backend/config"
	"github.com/LuckyPrima/attendance-backend/models"
	"github.com/gin-gonic/gin"
)

// ================== READ (All Employees, No Pagination) ==================
func GetEmployees(c *gin.Context) {
	var employees []models.Employee

	if err := config.DB.Preload("Department").Order("id ASC").Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

// ================== CREATE ==================
func CreateEmployee(c *gin.Context) {
	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// ================== UPDATE ==================
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee models.Employee

	if err := config.DB.First(&employee, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.Name = input.Name
	employee.Address = input.Address
	employee.DepartmentID = input.DepartmentID

	if err := config.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// ================== DELETE ==================
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Employee{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted"})
}
