package di

import (
	"fmt"
	"reflect"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/httputils"
	"github.com/gin-gonic/gin"
)

// New DI middleware
func New(app application.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		method := fmt.Sprintf("%v@%v", c.Request.Method, c.FullPath())
		runtimeKeyMap := application.DecodeHTTPRuntimeKey(c, app)
		tContext := app.ContextPool().Acquire(app, runtimeKeyMap, app.DB())
		if app.Conf().GetAtomicRequest() {
			tContext.GetTXDB()
		}
		defer func() {
			//release tcontext obj
			app.ContextPool().Release(tContext)
		}()
		controller, toFreeContainer := app.GetControllerPool().GetController(method, tContext, app, c)
		defer func() {
			app.GetControllerPool().Release(method, controller)
			for _, v := range toFreeContainer {
				app.GetContainerPool().Release(v)
			}
		}()
		controllerValue := reflect.ValueOf(controller) // new transport value
		controllerType := reflect.TypeOf(controller)   // transport type
		funcName, ok := app.GetControllerPool().GetControllerFuncName(method)
		if !ok {
			// if func not register , using the default method
			funcName = c.Request.Method
		}
		currentMethod, ok := controllerType.MethodByName(funcName)
		if !ok {
			panic("controller has no method ")
		}
		var inParam []reflect.Value                   // 构造函数入参 ， 入参1 ， transport指针对象 ， 入参2 ， context ， 入参3 ，pb  request
		inParam = append(inParam, controllerValue)    // 传入transport对象
		inParam = append(inParam, reflect.ValueOf(c)) // 传入ctx value
		// fmt.Println(currentMethod.Func.Type().NumIn())
		// to register controller in params
		// for i := 0; i < currentMethod.Func.Type().NumIn(); i++ {
		// 	t := currentMethod.Func.Type().In(i)
		// 	fmt.Println(t.Kind())
		// 	fmt.Println(t)
		// }
		res := currentMethod.Func.Call(inParam) // 调用transport函数，传入参数
		if len(res) != 3 {                      // 出参应该为2， 1为pb的response对象，2为error对象
			panic("wrong res type")
		}
		code, ok := res[0].Interface().(int)
		if !ok {
			panic("wrong code type")
		}
		if res[2].Interface() != nil {
			if app.Conf().GetAtomicRequest() {
				tContext.SafeRollback()
			}

			c.AbortWithStatusJSON(code, httputils.ResponseData{
				Status:  code,
				Result:  res[2].Interface().(error).Error(),
				Runtime: runtimeKeyMap,
			})
		} else {
			if app.Conf().GetAtomicRequest() {
				tContext.SafeCommit()
			}
			c.JSON(code, httputils.ResponseData{
				Status:  code,
				Result:  res[1].Interface(),
				Runtime: runtimeKeyMap,
			})
		}
	}
}
