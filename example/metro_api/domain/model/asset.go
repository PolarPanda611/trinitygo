package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// Asset model for Asset
type Asset struct {
	model.Model
	// to add your customize param inside here
	Code        string `json:"code" gorm:"type:varchar(50);index;not null;unique"`
	Name        string `json:"name" gorm:"type:varchar(50);"`
	Description string `json:"description" gorm:"type:varchar(50);"`
}
