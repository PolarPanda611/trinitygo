package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// User model for User
type User struct {
	model.Model
	// to add your customize param inside here
	UserName  string     `json:"user_name"  gorm:"type:varchar(50);unique"`
	NameLocal string     `json:"name_local"  gorm:"type:varchar(50)"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Languages []Language `json:"languages"`
}
