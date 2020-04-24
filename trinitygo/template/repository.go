package template

func init() {
	_templates["/domain/repository/repository.go"] = genRepository()
}

func genRepository() string {
	return `
package repository

// import (
// 	"fmt"
// 	"reflect"
// 	"sync"

// 	"github.com/PolarPanda611/trinitygo"
// 	"github.com/PolarPanda611/trinitygo/application"
// 	"{{.PackageName}}/domain/object"
// 	"{{.PackageName}}/domain/repository"
// )

// func init() {
// 	trinitygo.BindContainer(userRepositoryImpl{}}
// }

`
}
