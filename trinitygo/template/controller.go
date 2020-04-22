package template

func init() {
	_templates["/domain/controller/controller.go"] = genController()
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
// 	trinitygo.BindController("/YOURBASEURL",
// 		&sync.Pool{
// 			New: func() interface{} {
// 				controller := new("YOURCONTROLLER")
// 				return controller
// 			},
// 		},
// 		application.NewRequestMapping(httputil.GET, "YOURSUBPATH", "YOURFUNC"),
// 	)
// }
`
}
