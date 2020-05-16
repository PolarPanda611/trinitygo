package application

import (
	"github.com/PolarPanda611/trinitygo/conf"
	"github.com/PolarPanda611/trinitygo/keyword"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// AppMock for application mock
type AppMock struct {
	mock.Mock
}

// IsLogSelfCheck mock
func (a *AppMock) IsLogSelfCheck() bool {
	args := a.Called()
	return args.Get(0).(bool)
}

// Logger mock
func (a *AppMock) Logger() *golog.Logger {
	args := a.Called()
	return args.Get(0).(*golog.Logger)
}

// RuntimeKeys mock
func (a *AppMock) RuntimeKeys() []truntime.RuntimeKey {
	args := a.Called()
	return args.Get(0).([]truntime.RuntimeKey)
}

// Conf mock
func (a *AppMock) Conf() conf.Conf {
	args := a.Called()
	return args.Get(0).(conf.Conf)
}

// Keyword mock
func (a *AppMock) Keyword() keyword.Keyword {
	args := a.Called()
	return args.Get(0).(keyword.Keyword)
}

// ContextPool mock
func (a *AppMock) ContextPool() *ContextPool {
	args := a.Called()
	return args.Get(0).(*ContextPool)
}

// DB mock
func (a *AppMock) DB() *gorm.DB {
	args := a.Called()
	return args.Get(0).(*gorm.DB)
}

// Enforcer mock
func (a *AppMock) Enforcer() *casbin.Enforcer {
	args := a.Called()
	return args.Get(0).(*casbin.Enforcer)
}

// InstallDB mock
func (a *AppMock) InstallDB(f func() *gorm.DB) {}

// ControllerPool mock
func (a *AppMock) ControllerPool() *ControllerPool {
	args := a.Called()
	return args.Get(0).(*ControllerPool)
}

// InstancePool mock
func (a *AppMock) InstancePool() *InstancePool {
	args := a.Called()
	return args.Get(0).(*InstancePool)
}

// UseInterceptor mock
func (a *AppMock) UseInterceptor(interceptor ...grpc.UnaryServerInterceptor) Application {
	args := a.Called(interceptor)
	return args.Get(0).(Application)
}

// UseMiddleware mock
func (a *AppMock) UseMiddleware(middleware ...gin.HandlerFunc) Application {
	args := a.Called(middleware)
	return args.Get(0).(Application)
}

// RegRuntimeKey mock
func (a *AppMock) RegRuntimeKey(runtime ...truntime.RuntimeKey) Application {
	args := a.Called(runtime)
	return args.Get(0).(Application)
}

// InitGRPC mock
func (a *AppMock) InitGRPC() {}

// InitHTTP mock
func (a *AppMock) InitHTTP() {}

// InitRouter mock
func (a *AppMock) InitRouter() {}

// GetGRPCServer mock
func (a *AppMock) GetGRPCServer() *grpc.Server {
	args := a.Called()
	return args.Get(0).(*grpc.Server)
}

// ServeGRPC mock
func (a *AppMock) ServeGRPC() {}

// ServeHTTP mock
func (a *AppMock) ServeHTTP() {}

// ResponseFactory mock
func (a *AppMock) ResponseFactory() func(status int, res interface{}, runtime map[string]string) interface{} {
	args := a.Called()
	return args.Get(0).(func(status int, res interface{}, runtime map[string]string) interface{})
}
