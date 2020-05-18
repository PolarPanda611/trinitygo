package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// Group model for Group
type Group struct {
	model.Model
	// to add your customize param inside here
	Code string `json:"code" gorm:"type:varchar(50);index;not null;unique"`
	Name string `json:"name" gorm:"type:varchar(50);"`
}
