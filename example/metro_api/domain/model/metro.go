package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// Metro  列车信息
type Metro struct {
	model.Model
	Code        string `json:"code" gorm:"type:varchar(50);index;not null;"`
	Name        string `json:"name" gorm:"type:varchar(50);"`
	Description string `json:"description" gorm:"type:varchar(50);"`
}
