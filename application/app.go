package application

import (
	"context"

	truntime "github.com/PolarPanda611/trinitygo/runtime"

	"github.com/PolarPanda611/trinitygo/conf"

	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Application global app interface
type Application interface {
	Logger() *golog.Logger
	RuntimeKeys() []truntime.RuntimeKey
	Conf() conf.Conf
	ContextPool() *ContextPool
	DB() *gorm.DB
	InstallDB(f func() *gorm.DB)
	GetControllerPool() *ControllerPool
	GetServicePool() *ServicePool
	GetRepositoryPool() *RepositoryPool
	UseInterceptor(interceptor ...grpc.UnaryServerInterceptor) Application
	RegRuntimeKey(runtime ...truntime.RuntimeKey) Application
	InitGRPC()
	GetGRPCServer() *grpc.Server
	ServeGRPC()
}

// DecodeRuntimeKey  decode runtime key from ctx
func DecodeRuntimeKey(ctx context.Context, app Application) map[string]string {
	runtimeKeyMap := make(map[string]string)
	if ctx != nil {
		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			for _, v := range app.RuntimeKeys() {
				runtimeKeyMap[v.GetKeyName()] = md[v.GetKeyName()][0]
			}
		}
	}
	return runtimeKeyMap
}
