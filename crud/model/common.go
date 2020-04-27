package model

import (
	"time"
)

//Model common type
type Model struct {
	ID           int64      `json:"key,string"  gorm:"primary_key;AUTO_INCREMENT:false" mapstructure:"key"`
	CreatedTime  *time.Time `json:"created_time" `
	CreateUser   *User      `json:"create_user" `
	CreateUserID int64      `json:"create_user_id,string" `
	UpdatedTime  *time.Time `json:"updated_time"`
	UpdateUser   *User      `json:"update_user"`
	UpdateUserID int64      `json:"update_user_id,string" `
	DeletedTime  *time.Time `json:"deleted_time"`
	DeleteUser   *User      `json:"delete_user"`
	DeleteUserID int64      `json:"delete_user_id,string" `
	DVersion     string     `json:"d_version"`
}

//Simpmodel common type
type Simpmodel struct {
	ID          int64      `json:"key,string"  gorm:"primary_key;AUTO_INCREMENT:false"`
	CreatedTime *time.Time `json:"created_time" `
}

//Logmodel common type
type Logmodel struct {
	ID           int64      `json:"key,string"  gorm:"primary_key;AUTO_INCREMENT:false"`
	CreatedTime  *time.Time `json:"created_time" `
	CreateUser   *User      `json:"create_user"`
	CreateUserID int64      `json:"create_user_id,string" `
}

//Viewmodel for view type
type Viewmodel struct {
	ID           int64      `json:"key,string"  gorm:"primary_key;AUTO_INCREMENT:false"`
	CreatedTime  *time.Time `json:"created_time" `
	CreateUser   *User      `json:"create_user"`
	CreateUserID int64      `json:"create_user_id,string"`
	UpdatedTime  *time.Time `json:"updated_time" `
	UpdateUser   *User      `json:"update_user"`
	UpdateUserID int64      `json:"update_user_id,string"`
	DeletedTime  *time.Time `json:"deleted_time" `
	DeleteUser   *User      `json:"delete_user"`
	DeleteUserID int64      `json:"delete_user_id,string"`
}

// AutoModel with id auto increment
type AutoModel struct {
	ID           int64      `json:"key,string"  gorm:"primary_key;"`
	CreatedTime  *time.Time `json:"created_time" `
	CreateUser   *User      `json:"create_user"`
	CreateUserID int64      `json:"create_user_id,string"`
	UpdatedTime  *time.Time `json:"updated_time" `
	UpdateUser   *User      `json:"update_user"`
	UpdateUserID int64      `json:"update_user_id,string"`
	DeletedTime  *time.Time `json:"deleted_time" `
	DeleteUser   *User      `json:"delete_user"`
	DeleteUserID int64      `json:"delete_user_id,string"`
}
