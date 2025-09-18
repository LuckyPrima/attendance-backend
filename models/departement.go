package models

import "time"

type Department struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DepartementName string    `gorm:"size:255;not null" json:"departement_name"`
	MaxClockIn      string    `gorm:"column:max_clock_in;type:time" json:"max_clock_in"`
	MaxClockOut     string    `gorm:"column:max_clock_out;type:time" json:"max_clock_out"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Employees []Employee `gorm:"foreignKey:DepartmentID" json:"employees"`
}
