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
// 	"{{.PackageName}}/domain/object"
// 	"{{.PackageName}}/domain/repository"
// )

// func init() {
// 	trinitygo.BindContainer(userServiceImpl{})
// }

	`
}
