package model

//User model User
type User struct {
	Model
	UserName           string       `json:"user_name" gorm:"type:varchar(50);index;not null;"`            // login username /profile
	NameLocal          string       `json:"name_local"  gorm:"type:varchar(50);" `                        // local name
	NameEN             string       `json:"name_en"  gorm:"type:varchar(50);" `                           // EN name
	Email              string       `json:"email"  gorm:"type:varchar(50);" `                             // login email
	Phone              string       `json:"phone" gorm:"type:varchar(50);" `                              // login phone
	Groups             []Group      `json:"groups" gorm:"many2many:user_groups;"`                         // foreign key -->group
	Permissions        []Permission `json:"permissions" gorm:"many2many:user_permissions;"`               // foreign key --->permission
	PreferenceLanguage string       `json:"preference_language" gorm:"type:varchar(50);default:'en-US'" ` // user preference language

}
