package recovery

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/PolarPanda611/trinitygo/application"

	"google.golang.org/grpc"
)

// New recovery from panic
// initial Recovery , shoould be first interceptor
func New(app application.Application) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if err := recover(); err != nil {
				// when stack finishes
				logMessage := fmt.Sprintf("Recovered from GRPC Method %v \n", info.FullMethod)
				logMessage += fmt.Sprintf("At Request: %v\n", req)
				logMessage += fmt.Sprintf("Trace: %s\n", err)
				logMessage += fmt.Sprintf("\n%s", debug.Stack())
				app.Logger().Warn(logMessage)
			}
		}()
		return handler(ctx, req)
	}
}

// func Interceptor(app app.Application)
