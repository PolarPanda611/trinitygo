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
// 	"github.com/PolarPanda611/trinitygo/example/http/domain/object"
// 	"github.com/PolarPanda611/trinitygo/example/http/domain/repository"
// )

// func init() {
// 	trinitygo.BindContainer(reflect.TypeOf("&userRepositoryImpl{} <- YOURSERVICE"), &sync.Pool{
// 		New: func() interface{} {
// 			repository := new(userRepositoryImpl  < - YOURREPOSITORY)
// 			return repository
// 		},
// 	})
// }

	`
}
