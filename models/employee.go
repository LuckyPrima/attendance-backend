package models

import "time"

type Employee struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentID uint      `json:"department_id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Address      string    `gorm:"type:text" json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Department  Department   `gorm:"foreignKey:DepartmentID" json:"department"`
	Attendances []Attendance `gorm:"foreignKey:EmployeeID" json:"attendances"`
}
