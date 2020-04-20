package object

// Language model
type Language struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Code    string `json:"code" gorm:"type:varchar(50);unique"`
	Name    string `json:"name"  gorm:"type:varchar(50)"`
	UserID  uint   `json:"user_id"`
	UserXXX *User  `json:"user"`
}
