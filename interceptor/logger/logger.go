package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/PolarPanda611/trinitygo/logger"

	"github.com/PolarPanda611/trinitygo/application"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type interceptorLogger struct {
	app    application.Application
	config *logger.Config
}

// New record log
func New(app application.Application, config ...*logger.Config) grpc.UnaryServerInterceptor {
	c := logger.DefaultConfig()
	if len(config) > 0 {
		c = config[0]
	}

	l := &interceptorLogger{app: app, config: c}
	return l.Interceptor()
}

func (l *interceptorLogger) Interceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// skip the method
		for _, v := range l.config.Skippers {
			if v(info.FullMethod) {
				return handler(ctx, req)
			}
		}

		var method string
		var latency time.Duration
		var startTime, endTime time.Time
		startTime = time.Now()

		resp, err := handler(ctx, req)

		endTime = time.Now()

		if l.config.Method {
			method = info.FullMethod
		}
		if l.config.Latency {
			latency = endTime.Sub(startTime)
		}

		// record log

		// print the logs
		// if logFunc := l.config.LogFunc; logFunc != nil {
		// 	logFunc(endTime, latency, status, ip, method, path, message, headerMessage)
		// 	return
		// }

		line := fmt.Sprintf("%v %4v ", method, latency)
		if l.config.Runtime {
			for _, v := range l.app.RuntimeKeys() {
				md, _ := metadata.FromIncomingContext(ctx)
				line += fmt.Sprintf("%v %v ", v.GetKeyName(), md[v.GetKeyName()][0])
			}
		}
		if l.config.Request {
			line += fmt.Sprintf("%v %v ", "Request", req)
		}
		if l.config.Response {
			line += fmt.Sprintf("%v %v ", "Response", resp)
		}
		if l.config.Error && err != nil {
			line += fmt.Sprintf("%v %v ", "Error", err)
		}
		l.app.Logger().Info(line)
		return resp, err
	}
}
