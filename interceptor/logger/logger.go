package logger

import (
	"github.com/PolarPanda611/trinitygo/application"

	"google.golang.org/grpc"
)

// New record log
func New(app application.Application, config ...*application.LogConfig) grpc.UnaryServerInterceptor {
	c := application.DefaultLogConfig()
	if len(config) > 0 {
		c = config[0]
	}

	l := application.NewLogLogger(app, c)
	return l.Interceptor()
}
