package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// User model for User
type User struct {
	model.Model
	// to add your customize param inside here
	Code string `json:"code" gorm:"type:varchar(50);index;not null;unique"`
}
