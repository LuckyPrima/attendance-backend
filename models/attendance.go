package models

import "time"

type Attendance struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	EmployeeID uint       `json:"employee_id"`
	ClockIn    *time.Time `json:"clock_in"`
	ClockOut   *time.Time `json:"clock_out"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	Employee  Employee            `gorm:"foreignKey:EmployeeID;references:ID" json:"employee"`
	Histories []AttendanceHistory `gorm:"foreignKey:AttendanceID;references:ID" json:"histories"`
}
