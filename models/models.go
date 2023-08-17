package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	ID            string `gorm:"primaryKey"`
	Name          string `json:"name" gorm:"unique"`
	Description   string `json:"description"`
	EmployeeCount *int64 `json:"employee_count"`
	Registered    *bool  `json:"registered"`
	Type          string `json:"type"`
}

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}
