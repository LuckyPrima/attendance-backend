package routes

import (
	"github.com/LuckyPrima/attendance-backend/controllers"
	"github.com/gin-gonic/gin"
)

func AttendanceRoutes(r *gin.Engine) {
	att := r.Group("/attendance")
	{
		att.POST("/clockin", controllers.ClockIn)
		att.PUT("/clockout", controllers.ClockOut)
		att.GET("/logs", controllers.GetAttendanceLogs)
	}
}
