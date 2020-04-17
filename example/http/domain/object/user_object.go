package object

// User model
type User struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	UserName  string     `json:"user_name"  gorm:"type:varchar(50);unique"`
	NameLocal string     `json:"name_local"  gorm:"type:varchar(50)"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	Languages []Language `json:"languages"`
}
