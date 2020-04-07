package di

import (
	"context"
	"reflect"
	"strings"

	"github.com/PolarPanda611/trinitygo/application"

	"google.golang.org/grpc"
)

// New new DI interceptor
func New(app application.Application) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		runtimeKeyMap := application.DecodeRuntimeKey(ctx, app)
		tContext := app.ContextPool().Acquire(app, runtimeKeyMap, app.DB())
		method := strings.Split(info.FullMethod, "/") // /user.UserService/GetUserByID
		defer func() {
			//release tcontext obj
			app.ContextPool().Release(tContext)
		}()

		controller, toFreeService, toFreeRepo := app.GetControllerPool().GetController(method[1], tContext, app)
		defer func() {
			app.GetControllerPool().Release(method[1], controller)
			for _, v := range toFreeService {
				app.GetServicePool().Release(v)
			}
			for _, v := range toFreeRepo {
				app.GetRepositoryPool().Release(v)
			}
		}()
		controllerValue := reflect.ValueOf(controller) // new transport value
		controllerType := reflect.TypeOf(controller)   // transport type
		var currentMethod reflect.Method
		for i := 0; i < controllerType.NumMethod(); i++ { // 遍历transport all method
			m := controllerType.Method(i) //   get method
			if m.Name == method[2] {      // m.Name method name ,
				currentMethod = m
			}
		}
		var inParam []reflect.Value                     // 构造函数入参 ， 入参1 ， transport指针对象 ， 入参2 ， context ， 入参3 ，pb  request
		inParam = append(inParam, controllerValue)      // 传入transport对象
		inParam = append(inParam, reflect.ValueOf(ctx)) // 传入ctx value
		inParam = append(inParam, reflect.ValueOf(req)) // 传入pb request
		res := currentMethod.Func.Call(inParam)         // 调用transport函数，传入参数
		if len(res) != 2 {                              // 出参应该为2， 1为pb的response对象，2为error对象
			panic("wrong res type")
		}
		var resErr error
		if res[1].Interface() != nil {
			resErr = res[1].Interface().(error)
		} else {
			resErr = nil
		}
		return res[0].Interface(), resErr
	}
}
