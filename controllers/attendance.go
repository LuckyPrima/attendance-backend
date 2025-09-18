package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/LuckyPrima/attendance-backend/config"
	"github.com/LuckyPrima/attendance-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// helper: compare only time-of-day (hour,min,sec)
func timeOfDay(t time.Time) (h, m, s int) {
	return t.Hour(), t.Minute(), t.Second()
}

func isTimeLE(a time.Time, b time.Time) bool { // a <= b
	ha, ma, sa := timeOfDay(a)
	hb, mb, sb := timeOfDay(b)
	if ha != hb {
		return ha < hb
	}
	if ma != mb {
		return ma < mb
	}
	return sa <= sb
}

func isTimeGE(a time.Time, b time.Time) bool { // a >= b
	ha, ma, sa := timeOfDay(a)
	hb, mb, sb := timeOfDay(b)
	if ha != hb {
		return ha > hb
	}
	if ma != mb {
		return ma > mb
	}
	return sa >= sb
}

// parse "HH:mm" or "HH:mm:ss"
func parseDeptTime(val string) (time.Time, error) {
	if t, err := time.Parse("15:04", val); err == nil {
		return t, nil
	}
	return time.Parse("15:04:05", val)
}

// ================= CLOCK IN =================
func ClockIn(c *gin.Context) {
	var body struct {
		EmployeeID uint `json:"employee_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	// load employee with department
	var emp models.Employee
	if err := config.DB.Preload("Department").First(&emp, body.EmployeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if there's an active clock-in (today & not clocked out)
	var activeAtt models.Attendance
	err := config.DB.
		Where("employee_id = ? AND DATE(clock_in) = CURDATE() AND clock_out IS NULL", emp.ID).
		Order("id desc").
		First(&activeAtt).Error

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already clocked in, please clock out first"})
		return
	}

	// Create new attendance record
	newAtt := models.Attendance{
		EmployeeID: emp.ID,
		ClockIn:    &now,
	}
	if err := config.DB.Create(&newAtt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// status (On Time / Late)
	desc := "On Time"
	if deptIn, err := parseDeptTime(emp.Department.MaxClockIn); err == nil {
		if !isTimeLE(now, deptIn) {
			desc = "Late"
		}
	}

	// Save history
	history := models.AttendanceHistory{
		EmployeeID:     emp.ID,
		AttendanceID:   newAtt.ID,
		DateAttendance: now,
		AttendanceType: 1, // In
		Description:    desc,
	}
	if err := config.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "clock in recorded",
		"attendance": newAtt,
		"history":    history,
	})
}

// ================= CLOCK OUT =================
func ClockOut(c *gin.Context) {
	var body struct {
		EmployeeID uint `json:"employee_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	// load employee with department
	var emp models.Employee
	if err := config.DB.Preload("Department").First(&emp, body.EmployeeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if there's an active clock-in (today & not clocked out)
	var activeAtt models.Attendance
	err := config.DB.
		Where("employee_id = ? AND DATE(clock_in) = CURDATE() AND clock_out IS NULL", emp.ID).
		Order("id desc").
		First(&activeAtt).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no active clock-in found, please clock in first"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// update clock out
	activeAtt.ClockOut = &now
	if err := config.DB.Save(&activeAtt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// status (On Time / Early)
	desc := "On Time"
	if deptOut, err := parseDeptTime(emp.Department.MaxClockOut); err == nil {
		if !isTimeGE(now, deptOut) {
			desc = "Early"
		}
	}

	// Save history
	history := models.AttendanceHistory{
		EmployeeID:     emp.ID,
		AttendanceID:   activeAtt.ID,
		DateAttendance: now,
		AttendanceType: 2, // Out
		Description:    desc,
	}
	if err := config.DB.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "clock out recorded",
		"attendance": activeAtt,
		"history":    history,
	})
}

// ================= GET LOGS =================
func GetAttendanceLogs(c *gin.Context) {
	dateStr := c.Query("date")
	deptStr := c.Query("department_id")
	pageStr := c.Query("page")
	limit := 10
	page := 1

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	offset := (page - 1) * limit

	// Output struct
	type Out struct {
		HistoryID      uint       `json:"history_id"`
		EmployeeID     uint       `json:"employee_id"`
		EmployeeName   string     `json:"employee_name"`
		DepartmentID   uint       `json:"department_id"`
		DepartmentName string     `json:"department_name"`
		DateAttendance time.Time  `json:"date_attendance"`
		Type           string     `json:"type"`
		Description    string     `json:"description"`
		Accuracy       string     `json:"accuracy"`
	}

	var outs []Out

	// Base query
	db := config.DB.Table("attendance_histories as h").
		Select(`h.id as history_id, e.id as employee_id, e.name as employee_name,
                d.id as department_id, d.departement_name as department_name,
                h.date_attendance, 
                CASE h.attendance_type 
                    WHEN 1 THEN 'Clock In'
                    WHEN 2 THEN 'Clock Out'
                END as type,
                h.description`).
		Joins("join employees e on e.id = h.employee_id").
		Joins("join departments d on d.id = e.department_id")

	// Filter by date
	if dateStr != "" {
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (use YYYY-MM-DD)"})
			return
		}
		db = db.Where("DATE(h.date_attendance) = ?", dateStr)
	}

	// Filter by department
	if deptStr != "" {
		if _, err := strconv.Atoi(deptStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid department_id"})
			return
		}
		db = db.Where("e.department_id = ?", deptStr)
	}

	// Count total
	var total int64
	if err := db.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Execute query with pagination And sorting
	if err := db.
		Order("h.date_attendance DESC").
		Limit(limit).
		Offset(offset).
		Scan(&outs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Determine Accuracy (On Time / Late / Early)
	for i := range outs {
		var dept models.Department
		if err := config.DB.First(&dept, outs[i].DepartmentID).Error; err != nil {
			continue
		}

		attendanceTime := outs[i].DateAttendance

		if outs[i].Type == "Clock In" {
			if deptIn, err := parseDeptTime(dept.MaxClockIn); err == nil {
				if isTimeLE(attendanceTime, deptIn) {
					outs[i].Accuracy = "On Time"
				} else {
					outs[i].Accuracy = "Late"
				}
			}
		}

		if outs[i].Type == "Clock Out" {
			if deptOut, err := parseDeptTime(dept.MaxClockOut); err == nil {
				if isTimeGE(attendanceTime, deptOut) {
					outs[i].Accuracy = "On Time"
				} else {
					outs[i].Accuracy = "Early"
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       outs,
		"page":       page,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}