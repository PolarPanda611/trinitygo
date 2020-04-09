package di

import (
	"fmt"
	"reflect"

	"github.com/PolarPanda611/trinitygo/application"
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Status  int         // the http response status  to return
	Result  interface{} // the response data  if req success
	Runtime map[string]string
}

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
		controller, toFreeService, toFreeRepo := app.GetControllerPool().GetController(method, tContext, app)
		defer func() {
			app.GetControllerPool().Release(method, controller)
			for _, v := range toFreeService {
				app.GetServicePool().Release(v)
			}
			for _, v := range toFreeRepo {
				app.GetRepositoryPool().Release(v)
			}
		}()
		controllerValue := reflect.ValueOf(controller) // new transport value
		controllerType := reflect.TypeOf(controller)   // transport type
		currentMethod, ok := controllerType.MethodByName(c.Request.Method)
		if !ok {
			panic("controller has no method ")
		}
		var inParam []reflect.Value                   // 构造函数入参 ， 入参1 ， transport指针对象 ， 入参2 ， context ， 入参3 ，pb  request
		inParam = append(inParam, controllerValue)    // 传入transport对象
		inParam = append(inParam, reflect.ValueOf(c)) // 传入ctx value
		res := currentMethod.Func.Call(inParam)       // 调用transport函数，传入参数
		if len(res) != 3 {                            // 出参应该为2， 1为pb的response对象，2为error对象
			panic("wrong res type")
		}
		code, ok := res[0].Interface().(int)
		if !ok {
			panic("wrong code type")
		}
		if res[2].Interface() != nil {
			if app.Conf().GetAtomicRequest() {
				tContext.GetDB().Rollback()
			}

			c.AbortWithStatusJSON(code, ResponseData{
				Status:  code,
				Result:  res[2].Interface().(error).Error(),
				Runtime: runtimeKeyMap,
			})
		} else {
			if app.Conf().GetAtomicRequest() {
				tContext.GetDB().Commit()
			}
			c.JSON(code, ResponseData{
				Status:  code,
				Result:  res[1].Interface(),
				Runtime: runtimeKeyMap,
			})
		}
	}
}
