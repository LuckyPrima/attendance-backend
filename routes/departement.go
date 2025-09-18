package routes

import (
	"github.com/LuckyPrima/attendance-backend/controllers"
	"github.com/gin-gonic/gin"
)

func DepartmentRoutes(r *gin.Engine) {
	r.GET("/departments", controllers.GetDepartments)
	r.POST("/departments", controllers.CreateDepartment)
	r.PUT("/departments/:id", controllers.UpdateDepartment)
	r.DELETE("/departments/:id", controllers.DeleteDepartment)
}
