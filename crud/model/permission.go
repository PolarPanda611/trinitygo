package model

//Permission model Role
type Permission struct {
	Model
	Code string `json:"code" gorm:"type:varchar(100);index;unique;not null;"`
	Name string `json:"name" gorm:"type:varchar(100);not null;default:''"`
}
