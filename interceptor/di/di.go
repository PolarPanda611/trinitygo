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
		runtimeKeyMap := application.DecodeGRPCRuntimeKey(ctx, app.RuntimeKeys())
		tContext := app.ContextPool().Acquire(app, runtimeKeyMap, app.DB(), nil)
		if app.Conf().GetAtomicRequest() {
			tContext.DBTx()
		}
		method := strings.Split(info.FullMethod, "/") // /user.UserService/GetUserByID
		defer func() {
			//release tcontext obj
			app.ContextPool().Release(tContext)
		}()

		controller, sharedInstance := app.ControllerPool().GetController(method[1], tContext, app, nil)
		defer func() {
			app.ControllerPool().Release(method[1], controller)
			for _, v := range sharedInstance {
				app.ContainerPool().Release(v)
			}
		}()
		currentMethod, ok := reflect.TypeOf(controller).MethodByName(method[2])
		if !ok {
			panic("controller has no method ")
		}
		var inParam []reflect.Value
		inParam = append(inParam, reflect.ValueOf(controller))
		inParam = append(inParam, reflect.ValueOf(ctx))
		inParam = append(inParam, reflect.ValueOf(req))
		res := currentMethod.Func.Call(inParam)
		if len(res) != 2 {
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
