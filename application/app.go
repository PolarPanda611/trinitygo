package application

import (
	"context"

	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/gin-gonic/gin"

	"github.com/PolarPanda611/trinitygo/conf"

	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Application global app interface
type Application interface {
	IsLogSelfCheck() bool
	Logger() *golog.Logger
	RuntimeKeys() []truntime.RuntimeKey
	Conf() conf.Conf
	ContextPool() *ContextPool
	DB() *gorm.DB
	InstallDB(f func() *gorm.DB)
	ControllerPool() *ControllerPool
	InstancePool() *InstancePool
	UseInterceptor(interceptor ...grpc.UnaryServerInterceptor) Application
	UseMiddleware(middleware ...gin.HandlerFunc) Application
	RegRuntimeKey(runtime ...truntime.RuntimeKey) Application
	InitGRPC()
	InitHTTP()
	InitRouter()
	GetGRPCServer() *grpc.Server
	ServeGRPC()
	ServeHTTP()
}

// DecodeGRPCRuntimeKey  decode runtime key from ctx
func DecodeGRPCRuntimeKey(ctx context.Context, runtimeKey []truntime.RuntimeKey) map[string]string {
	runtimeKeyMap := make(map[string]string)
	if ctx != nil {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			for _, v := range runtimeKey {
				runtimeKeyMap[v.GetKeyName()] = md[v.GetKeyName()][0]
			}
		}
	}
	return runtimeKeyMap
}

// DecodeHTTPRuntimeKey decode http runtime
func DecodeHTTPRuntimeKey(c *gin.Context, runtimeKey []truntime.RuntimeKey) map[string]string {
	runtimeKeyMap := make(map[string]string)
	if c != nil {
		for _, v := range runtimeKey {
			runtimeKeyMap[v.GetKeyName()] = c.GetString(v.GetKeyName())
		}
	}
	return runtimeKeyMap
}
