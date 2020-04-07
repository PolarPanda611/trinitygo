package trinitygo

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	truntime "github.com/PolarPanda611/trinitygo/runtime"

	"github.com/PolarPanda611/trinitygo/interceptor/logger"

	"github.com/PolarPanda611/trinitygo/interceptor/di"

	"github.com/PolarPanda611/trinitygo/db"

	"github.com/PolarPanda611/trinitygo/conf"

	"github.com/PolarPanda611/trinitygo/interceptor/recovery"
	"github.com/PolarPanda611/trinitygo/sd"
	"github.com/PolarPanda611/trinitygo/utils"

	"github.com/PolarPanda611/trinitygo/interceptor/runtime"

	"github.com/PolarPanda611/trinitygo/application"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
	"google.golang.org/grpc"
)

var (
	configpath     string = "./config/"
	controllerPool *application.ControllerPool
	servicePool    *application.ServicePool
	repositoryPool *application.RepositoryPool
)

func init() {
	controllerPool = application.NewControllerPool()
	servicePool = application.NewServicePool()
	repositoryPool = application.NewRepositoryPool()
}

// Application core of trinity
type Application struct {
	config      conf.Conf
	logger      *golog.Logger
	contextPool *application.ContextPool

	// used for build
	once           sync.Once
	mu             sync.RWMutex
	db             *gorm.DB
	controllerPool *application.ControllerPool
	servicePool    *application.ServicePool
	repositoryPool *application.RepositoryPool
	interceptors   []grpc.UnaryServerInterceptor
	runtimeKeys    []truntime.RuntimeKey
	serviceMesh    sd.ServiceMesh

	grpcServer *grpc.Server
}

// SetConfigPath set config path
func SetConfigPath(path string) {
	configpath = path
}

// New new application
func New() application.Application {
	app := &Application{
		logger: golog.Default,
	}
	app.config = conf.NewSetting(configpath)

	appPrefix := fmt.Sprintf("[%v@%v]", utils.GetServiceName(app.config.GetProjectName()), app.config.GetProjectVersion())
	app.logger.SetPrefix(appPrefix)
	app.logger.SetTimeFormat("2006-01-02 15:04:05.000")

	app.contextPool = application.New(func() application.Context {
		return application.NewContext(app)
	})

	app.controllerPool = controllerPool
	app.servicePool = servicePool
	app.repositoryPool = repositoryPool
	return app
}

// BindController bind service
func BindController(controllerName string, Pool *sync.Pool) {
	controllerPool.NewController(controllerName, Pool)
}

// BindService bind service
func BindService(serviceName reflect.Type, Pool *sync.Pool) {
	servicePool.NewService(serviceName, Pool)
}

// BindRepository bind service
func BindRepository(repoName reflect.Type, Pool *sync.Pool) {
	repositoryPool.NewRepository(repoName, Pool)
}

// DefaultGRPC default grpc server
func DefaultGRPC() application.Application {
	app := New()
	app.UseInterceptor(recovery.New(app))
	app.UseInterceptor(runtime.New(app))
	app.UseInterceptor(logger.New(app))
	app.UseInterceptor(di.New(app))
	// app.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", true, ""))
	// app.RegRuntimeKey(truntime.NewRuntimeKey("user_id", true, ""))
	// app.RegRuntimeKey(truntime.NewRuntimeKey("user_name", true, ""))
	return app
}

// Logger get logger
func (app *Application) Logger() *golog.Logger {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.logger
}

// RuntimeKeys get runtime keys
func (app *Application) RuntimeKeys() []truntime.RuntimeKey {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.runtimeKeys
}

// Conf get conf
func (app *Application) Conf() conf.Conf {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.config
}

// DB get db
func (app *Application) DB() *gorm.DB {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.db
}

// ContextPool get contextpoo;
func (app *Application) ContextPool() *application.ContextPool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.contextPool
}

// GetControllerPool get all serviice pool
func (app *Application) GetControllerPool() *application.ControllerPool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.controllerPool
}

// GetServicePool get all serviice pool
func (app *Application) GetServicePool() *application.ServicePool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.servicePool
}

// GetRepositoryPool get all serviice pool
func (app *Application) GetRepositoryPool() *application.RepositoryPool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.repositoryPool
}

// UseInterceptor application use interceptor
func (app *Application) UseInterceptor(interceptor ...grpc.UnaryServerInterceptor) application.Application {
	app.interceptors = append(app.interceptors, interceptor...)
	return app
}

// RegRuntimeKey register runtime key
func (app *Application) RegRuntimeKey(runtime ...truntime.RuntimeKey) application.Application {
	app.runtimeKeys = append(app.runtimeKeys, runtime...)
	return app
}

// InstallDB install db
func (app *Application) InstallDB(f func() *gorm.DB) {
	app.db = f()
	err := app.db.DB().Ping()
	if err != nil {
		app.Logger().Fatal("booting detected db initializer ...install failed ,err : %v", err)
	}
	app.Logger().Info("booting detected db initializer ...install successfully , test passed ! ")
}

func (app *Application) initPool() {
	poolKey := ""
	for _, v := range app.controllerPool.GetControllerMap() {
		poolKey += fmt.Sprintf("%v,", v)
	}
	// service pool checking
	line := fmt.Sprintf("booting detected %v controller pool (%v)...installed", len(app.controllerPool.GetControllerMap()), poolKey)
	app.Logger().Info(line)

	poolKey = ""
	for _, v := range app.servicePool.GetServiceType() {
		poolKey += fmt.Sprintf("%v,", v)
	}
	// service pool checking
	line = fmt.Sprintf("booting detected %v service pool (%v)...installed", len(app.servicePool.GetServiceType()), poolKey)
	app.Logger().Info(line)

	poolKey = ""
	for _, v := range app.repositoryPool.GetRepositoryType() {
		poolKey += fmt.Sprintf("%v,", v)
	}
	// service pool checking
	line = fmt.Sprintf("booting detected %v repository pool (%v)...installed", len(app.repositoryPool.GetRepositoryType()), poolKey)
	app.Logger().Info(line)
}

func (app *Application) initRuntime() {
	runtimeKey := ""
	for _, v := range app.runtimeKeys {
		runtimeKey += fmt.Sprintf("%v,", v.GetKeyName())
	}
	// runtime checking
	line := fmt.Sprintf("booting detected %v runtime([%v]) ...installed", len(app.runtimeKeys), runtimeKey)
	app.Logger().Info(line)
}
func (app *Application) initDB() {
	if app.db == nil {
		f := func() *gorm.DB {
			return db.DefaultInstallGORM(
				app.config.GetDebug(),
				true,
				app.config.GetDBType(),
				app.config.GetDBTablePrefix(),
				app.config.GetDBServer(),
				app.config.GetDbMaxIdleConn(),
				app.config.GetDbMaxOpenConn(),
			)
		}
		app.db = f()
		line := fmt.Sprintf("booting db instance with default install")
		app.Logger().Info(line)
	}
}

func (app *Application) initSD() {
	if app.config.GetServiceDiscoveryAutoRegister() {
		switch app.config.GetServiceDiscoveryType() {
		case "etcd":
			c, err := sd.NewEtcdRegister(
				app.config.GetServiceDiscoveryAddress(),
				app.config.GetServiceDiscoveryPort(),
			)
			if err != nil {
				app.logger.Fatal("get service mesh client err")
			}
			app.serviceMesh = c
			break
		default:
			app.logger.Fatal("wrong service mash type")
		}
	}
}

// InitTrinity serve grpc
func (app *Application) InitTrinity() {
	app.initRuntime()
	app.initDB()
	app.initPool()
	app.initSD()
}

// InitGRPC serve grpc
func (app *Application) InitGRPC() {
	app.InitTrinity()

	// interceptors installation
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			app.interceptors...,
		),
	}
	line := fmt.Sprintf("booting detected %v interceptors ...installed ", len(app.interceptors))
	app.Logger().Info(line)

	app.grpcServer = grpc.NewServer(opts...)
}

// GetGRPCServer get grpc server instance
func (app *Application) GetGRPCServer() *grpc.Server {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.grpcServer
}

// ServeGRPC serve grpc server
func (app *Application) ServeGRPC() {
	addr := fmt.Sprintf("%v:%v", app.config.GetAppAddress(), app.config.GetAppPort())
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("tcp port : %v  listen err: %v", addr, err)
	}
	gErr := make(chan error)
	go func() {
		if app.config.GetServiceDiscoveryAutoRegister() {
			if err := app.serviceMesh.RegService(
				app.config.GetProjectName(),
				app.config.GetProjectVersion(),
				app.config.GetAppAddress(),
				app.config.GetAppPort(),
				app.config.GetProjectTags(),
			); err != nil {
				gErr <- err
			}
		}
		line := fmt.Sprintf("booted grpc service listen at %v started", addr)
		app.Logger().Info(line)
		gErr <- app.grpcServer.Serve(lis)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		gErr <- fmt.Errorf("%s", <-c)
	}()

	line := fmt.Sprintf("booted grpc service listen at %v ,terminated err: %v ", addr, <-gErr)
	app.Logger().Info(line)
	if app.config.GetServiceDiscoveryAutoRegister() {
		app.serviceMesh.DeRegService(
			app.config.GetProjectName(),
			app.config.GetProjectVersion(),
			app.config.GetAppAddress(),
			app.config.GetAppPort(),
		)
	}
}
