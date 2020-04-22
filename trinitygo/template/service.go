package template

func init() {
	_templates["/domain/service/service.go"] = genService()
}

func genService() string {
	return `
package service

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
// 	trinitygo.BindContainer(reflect.TypeOf("&userServiceImpl{} <- YOURSERVICE"), &sync.Pool{
// 		New: func() interface{} {
// 			service := new(userServiceImpl  < - userServiceImpl)
// 			return service
// 		},
// 	})
// }

	`
}
