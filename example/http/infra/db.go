package infra

import (
	"github.com/PolarPanda611/trinitygo/example/http/domain/model"
	"github.com/jinzhu/gorm"
)

// DB app global DB instance
var DB *gorm.DB

// Migrate migrate table
func Migrate() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Language{})
	// userStruct := DB.NewScope(&object.User{}).GetModelStruct()
	// languageStruct := DB.NewScope(&object.Language{}).GetModelStruct()
	// fmt.Println(userStruct)
	// fmt.Println(languageStruct.StructFields[4].Struct.Type)
	// fmt.Println(languageStruct)
}
