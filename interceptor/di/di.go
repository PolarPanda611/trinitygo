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
		runtimeKeyMap := application.DecodeGRPCRuntimeKey(ctx, app)
		tContext := app.ContextPool().Acquire(app, runtimeKeyMap, app.DB())
		if app.Conf().GetAtomicRequest() {
			tContext.GetTXDB()
		}
		method := strings.Split(info.FullMethod, "/") // /user.UserService/GetUserByID
		defer func() {
			//release tcontext obj
			app.ContextPool().Release(tContext)
		}()

		controller, toFreeContainer := app.GetControllerPool().GetController(method[1], tContext, app, nil)
		defer func() {
			app.GetControllerPool().Release(method[1], controller)
			for _, v := range toFreeContainer {
				app.GetContainerPool().Release(v)
			}
		}()
		controllerValue := reflect.ValueOf(controller) // new transport value
		controllerType := reflect.TypeOf(controller)   // transport type
		currentMethod, ok := controllerType.MethodByName(method[2])
		if !ok {
			panic("controller has no method ")
		}
		var inParam []reflect.Value                     // 构造函数入参 ， 入参1 ， transport指针对象 ， 入参2 ， context ， 入参3 ，pb  request
		inParam = append(inParam, controllerValue)      // 传入transport对象
		inParam = append(inParam, reflect.ValueOf(ctx)) // 传入ctx value
		inParam = append(inParam, reflect.ValueOf(req)) // 传入pb request
		res := currentMethod.Func.Call(inParam)         // 调用transport函数，传入参数
		if len(res) != 2 {                              // 出参应该为2， 1为pb的response对象，2为error对象
			panic("wrong res type")
		}
		if res[1].Interface() != nil {
			tContext.SafeRollback()
			return nil, res[1].Interface().(error)
		}
		if app.Conf().GetAtomicRequest() {
			tContext.SafeCommit()
		}
		return res[0].Interface(), nil
	}
}
