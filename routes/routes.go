package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	EmployeeRoutes(r)
	DepartmentRoutes(r)
	AttendanceRoutes(r)
}
