package trinitygo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/gin-gonic/gin"

	"github.com/PolarPanda611/trinitygo/interceptor/logger"

	mlogger "github.com/PolarPanda611/trinitygo/middleware/logger"

	"github.com/PolarPanda611/trinitygo/interceptor/di"

	"github.com/PolarPanda611/trinitygo/db"

	"github.com/PolarPanda611/trinitygo/conf"

	"github.com/PolarPanda611/trinitygo/interceptor/recovery"
	"github.com/PolarPanda611/trinitygo/sd"
	"github.com/PolarPanda611/trinitygo/utils"

	"github.com/PolarPanda611/trinitygo/interceptor/runtime"
	mdi "github.com/PolarPanda611/trinitygo/middleware/di"
	httprecovery "github.com/PolarPanda611/trinitygo/middleware/recovery"
	mruntime "github.com/PolarPanda611/trinitygo/middleware/runtime"

	"github.com/PolarPanda611/trinitygo/application"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
	"google.golang.org/grpc"
)

var (
	configpath     string = "./config/"
	controllerPool *application.ControllerPool
	containerPool  *application.ContainerPool
)

func init() {
	controllerPool = application.NewControllerPool()
	containerPool = application.NewContainerPool()
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
	containerPool  *application.ContainerPool

	//grpc
	interceptors []grpc.UnaryServerInterceptor
	runtimeKeys  []truntime.RuntimeKey
	serviceMesh  sd.ServiceMesh
	grpcServer   *grpc.Server

	//http
	middlewares []gin.HandlerFunc
	router      *gin.Engine
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
	app.containerPool = containerPool
	return app
}

// BindController bind service
func BindController(controllerName string, Pool *sync.Pool) {
	controllerPool.NewController(controllerName, Pool)
}

// BindContainer bind container
func BindContainer(container reflect.Type, Pool *sync.Pool) {
	containerPool.NewContainer(container, Pool)
}

// DefaultGRPC default grpc server
func DefaultGRPC() application.Application {
	app := New()
	app.UseInterceptor(recovery.New(app))
	app.UseInterceptor(runtime.New(app))
	app.UseInterceptor(logger.New(app))
	app.UseInterceptor(di.New(app))
	// app.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", true,  func() string { return "" })
	// app.RegRuntimeKey(truntime.NewRuntimeKey("user_id", true,  func() string { return "" })
	// app.RegRuntimeKey(truntime.NewRuntimeKey("user_name", true,  func() string { return "" } )
	return app
}

// DefaultHTTP default http server
func DefaultHTTP() application.Application {
	app := New()
	app.UseMiddleware(mlogger.New(app))
	app.UseMiddleware(httprecovery.New(app))
	app.UseMiddleware(mruntime.New(app))

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

// GetContainerPool get all serviice pool
func (app *Application) GetContainerPool() *application.ContainerPool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.containerPool
}

// UseInterceptor application use interceptor, only impact on http server
func (app *Application) UseInterceptor(interceptor ...grpc.UnaryServerInterceptor) application.Application {
	app.interceptors = append(app.interceptors, interceptor...)
	return app
}

// UseMiddleware application use middleware , only impact on http server
func (app *Application) UseMiddleware(middleware ...gin.HandlerFunc) application.Application {
	app.middlewares = append(app.middlewares, middleware...)
	return app

}

// RegRuntimeKey register runtime key
// the runtime key should be lower case , because when the metadata transfer , it will all transform
// to lower case
func (app *Application) RegRuntimeKey(runtime ...truntime.RuntimeKey) application.Application {
	app.runtimeKeys = append(app.runtimeKeys, runtime...)
	return app
}

// InstallDB install db
func (app *Application) InstallDB(f func() *gorm.DB) {
	app.db = f()
	err := app.db.DB().Ping()
	if err != nil {
		app.Logger().Fatal("booting detected db initializer ...install failed ,err : ", err)
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
	for _, v := range app.containerPool.GetContainerType() {
		poolKey += fmt.Sprintf("%v,", v)
	}
	// service pool checking
	line = fmt.Sprintf("booting detected %v container pool (%v)...installed", len(containerPool.GetContainerType()), poolKey)
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

// InitHTTP serve grpc
func (app *Application) InitHTTP() {
	app.InitTrinity()
	gin.DefaultWriter = ioutil.Discard
	app.router = gin.New()
	for _, v := range app.middlewares {
		app.router.Use(v)
	}
	for _, controllerName := range app.GetControllerPool().GetControllerMap() {
		controllerNameList := strings.Split(controllerName, "@")
		method := controllerNameList[0]
		path := controllerNameList[1]
		switch method {
		case "GET":
			app.router.GET(path, mdi.New(app))
			break
		case "PATCH":
			app.router.PATCH(path, mdi.New(app))
			break
		case "POST":
			app.router.POST(path, mdi.New(app))
			break
		case "DELETE":
			app.router.DELETE(path, mdi.New(app))
			break
		default:
			app.logger.Fatal("booting err : wrong method ", method, " when binding controller ", controllerName)
		}
	}
	// app.router
}

// InitHTTP serve grpc
func (app *Application) ServeHTTP() {
	addr := fmt.Sprintf(":%v", app.config.GetAppPort())
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
			line := fmt.Sprintf("boooting http service registered successfully !")
			app.Logger().Info(line)
		}
		s := &http.Server{
			Addr:              addr,
			Handler:           app.router,
			ReadTimeout:       time.Duration(app.config.GetAppReadTimeout()) * time.Second,
			ReadHeaderTimeout: time.Duration(app.config.GetAppReadHeaderTimeout()) * time.Second,
			WriteTimeout:      time.Duration(app.config.GetAppWriteTimeout()) * time.Second,
			IdleTimeout:       time.Duration(app.config.GetAppIdleTimeout()) * time.Second,
			MaxHeaderBytes:    app.config.GetAppMaxHeaderBytes(),
		}
		line := fmt.Sprintf("booted http service listen at %v started ", addr)
		app.Logger().Info(line)
		gErr <- s.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		gErr <- fmt.Errorf("%s", <-c)
	}()

	line := fmt.Sprintf("booted http service listen at %v ,terminated err: %v ", addr, <-gErr)
	app.Logger().Error(line)
	app.db.Close()
	if app.config.GetServiceDiscoveryAutoRegister() {
		err := app.serviceMesh.DeRegService(
			app.config.GetProjectName(),
			app.config.GetProjectVersion(),
			app.config.GetAppAddress(),
			app.config.GetAppPort(),
		)
		if err != nil {
			line = fmt.Sprintf("boooting http service deregistered failed !")
			app.Logger().Error(line)
		} else {
			line = fmt.Sprintf("boooting http service deregistered successfully !")
			app.Logger().Info(line)
		}

	}
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
			line := fmt.Sprintf("boooting grpc service registered successfully !")
			app.Logger().Info(line)
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
	app.Logger().Error(line)
	app.db.Close()
	if app.config.GetServiceDiscoveryAutoRegister() {
		err := app.serviceMesh.DeRegService(
			app.config.GetProjectName(),
			app.config.GetProjectVersion(),
			app.config.GetAppAddress(),
			app.config.GetAppPort(),
		)
		if err != nil {
			line = fmt.Sprintf("boooting grpc service deregistered failed !")
			app.Logger().Error(line)
		} else {
			line = fmt.Sprintf("boooting grpc service deregistered successfully !")
			app.Logger().Info(line)
		}

	}

}
