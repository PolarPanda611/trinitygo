package model

import "github.com/PolarPanda611/trinitygo/crud/model"

// Line  metro line number
// 地铁线路
type Line struct {
	model.Model
	Code        string `json:"code" gorm:"type:varchar(50);index;not null;"`
	Name        string `json:"name" gorm:"type:varchar(50);"`
	Description string `json:"description" gorm:"type:varchar(50);"`
}
