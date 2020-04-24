package template

func init() {
	_templates["/domain/controller/http/controller.go"] = genController()
}

func genController() string {
	return `
package http

// import (
// 	"sync"

// 	"github.com/PolarPanda611/trinitygo"
// 	"github.com/PolarPanda611/trinitygo/application"
// 	"github.com/PolarPanda611/trinitygo/httputil"
// )

// func init() {
// 	trinitygo.BindController("/users", userControllerImpl{},
// 		application.NewRequestMapping(httputil.GET, "/:id", "GET"),
// 		application.NewRequestMapping(httputil.GET, "", "Getsssss"),
// 	)
// }
`
}
