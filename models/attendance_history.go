package models

import "time"

type AttendanceHistory struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeID     uint      `json:"employee_id"`
	AttendanceID   uint      `json:"attendance_id"`
	DateAttendance time.Time `json:"date_attendance"`
	AttendanceType int       `gorm:"type:tinyint" json:"attendance_type"` // 1 = In, 2 = Out
	Description    string    `gorm:"type:text" json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Employee   Employee   `gorm:"foreignKey:EmployeeID;references:ID" json:"employee"`
	Attendance Attendance `gorm:"foreignKey:AttendanceID;references:ID" json:"attendance"`
}
