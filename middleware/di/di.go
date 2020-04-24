package di

import (
	"fmt"
	"reflect"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/gin-gonic/gin"
)

// New DI middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := fmt.Sprintf("%v@%v", c.Request.Method, c.FullPath())
		runtimeKeyMap := application.DecodeHTTPRuntimeKey(c, app.RuntimeKeys())
		tContext := app.ContextPool().Acquire(app, runtimeKeyMap, app.DB(), c)
		defer func() {
			//release tcontext obj
			app.ContextPool().Release(tContext)
		}()
		controller, sharedInstance := app.ControllerPool().GetController(method, tContext, app, c)
		defer func() {
			for _, v := range sharedInstance {
				app.InstancePool().Release(v)
			}
		}()
		validators := app.ControllerPool().GetControllerValidators(method)
		for _, v := range validators {
			v(tContext)
			if c.IsAborted() {
				return
			}
		}
		funcName, ok := app.ControllerPool().GetControllerFuncName(method)
		if ok && funcName == "" || !ok {
			funcName = c.Request.Method
		}
		currentMethod, ok := reflect.TypeOf(controller).MethodByName(funcName)
		if !ok {
			panic("controller has no method ")
		}
		// validation passed , run handler
		var inParam []reflect.Value                            // 构造函数入参 ， 入参1 ， transport指针对象
		inParam = append(inParam, reflect.ValueOf(controller)) // 传入transport对象
		currentMethod.Func.Call(inParam)                       // 调用transport函数，传入参数
	}
}
