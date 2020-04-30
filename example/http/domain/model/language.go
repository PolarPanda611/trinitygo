package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// Language model for Language
type Language struct {
	model.Model
	// to add your customize param inside here
	Code string `json:"code" gorm:"type:varchar(50);unique"`
	Name string `json:"name"  gorm:"type:varchar(50)"`
}
