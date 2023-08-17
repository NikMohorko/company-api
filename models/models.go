package models

type Company struct {
	Id            string `gorm:"primaryKey;column:id"`
	Name          string `json:"name" gorm:"unique;column:name"`
	Description   string `json:"description" gorm:"column:description"`
	EmployeeCount *int64 `json:"employee_count" gorm:"column:employee_count"`
	Registered    *bool  `json:"registered" gorm:"column:registered"`
	Type          string `json:"type" gorm:"column:type"`
}

type User struct {
	Id       string `gorm:"primaryKey;column:id"`
	Username string `json:"username" gorm:"unique;column:username"`
	Password string `json:"password" gorm:"column:password"`
}
