package infra

import (
	"github.com/PolarPanda611/trinitygo/example/http/domain/object"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func Migrate() {
	DB.AutoMigrate(&object.User{})
	DB.AutoMigrate(&object.Language{})
}
