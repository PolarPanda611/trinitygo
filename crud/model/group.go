package model

//Group model Group
type Group struct {
	Model
	Name        string       `json:"name" gorm:"type:varchar(50);unique;not null;"`
	Description string       `json:"description" gorm:"type:varchar(100);not null;default:''"`
	Permissions []Permission `json:"permissions" gorm:"many2many:group_permissions;"`
	PGroup      *Group       `json:"p_group"`
	PID         int64        `json:"p_id,string"`
	Roles       []Role       `json:"roles" gorm:"many2many:group_roles;"`
}
