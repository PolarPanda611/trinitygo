package model

//Role model Role
type Role struct {
	Model
	Name        string       `json:"name" gorm:"type:varchar(50);index;unique;not null;"`
	Description string       `json:"description" gorm:"type:varchar(100);not null;default:''"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}
