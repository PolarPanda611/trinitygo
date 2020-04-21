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

	"github.com/PolarPanda611/trinitygo/httputil"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/PolarPanda611/trinitygo/interceptor/logger"

	"github.com/PolarPanda611/trinitygo/interceptor/di"
	mlogger "github.com/PolarPanda611/trinitygo/middleware/logger"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/PolarPanda611/trinitygo/db"

	"github.com/PolarPanda611/trinitygo/conf"

	"github.com/PolarPanda611/trinitygo/interceptor/recovery"
	"github.com/PolarPanda611/trinitygo/sd"

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
	_                          application.Application = new(Application)
	_app                       *Application
	_initApp                   sync.Once
	_configPath                string = "./config/"
	_bootingControllers        []bootingController
	_bootingContainers         []bootingContainer
	_healthCheckPath           string                                                 = "/ping"
	_enableHealthCheckPath     bool                                                   = false
	_defaultHealthCheckHandler func(app application.Application) func(c *gin.Context) = func(app application.Application) func(c *gin.Context) {
		return func(c *gin.Context) {
			err := app.DB().DB().Ping()
			if err != nil {
				c.AbortWithStatusJSON(400, httputil.ResponseData{
					Status: 400,
					Error:  err,
				})
				return
			}
			c.JSON(200, httputil.ResponseData{
				Status: 200,
				Result: gin.H{
					"APIStatus": "alive",
					"DBStatus":  "alive",
					"DBInfo":    app.DB().DB().Stats(),
				},
			})
			return
		}
	}
)

type bootingController struct {
	controllerName string
	controllerPool *sync.Pool
	requestMaps    []*application.RequestMap
}

type bootingContainer struct {
	containerName reflect.Type
	containerPool *sync.Pool
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
	serviceMesh  sd.ServiceMesh
	grpcServer   *grpc.Server

	//http
	middlewares []gin.HandlerFunc
	router      *gin.Engine

	//runtime
	runtimeKeys []truntime.RuntimeKey
}

// SetConfigPath set config path
func SetConfigPath(path string) {
	_configPath = path
}

// EnableHealthCheckURL set config path
func EnableHealthCheckURL(path ...string) {
	if len(path) > 0 {
		_healthCheckPath = path[0]
	}
	_enableHealthCheckPath = true
}

// SetHealthCheckDefaultHandler set default health check handler
func SetHealthCheckDefaultHandler(handler func(app application.Application) func(c *gin.Context)) {
	_defaultHealthCheckHandler = handler
}

// SetDefaultHeaderPrefix set default header prefix
func SetDefaultHeaderPrefix(newPrefix string) {
	mruntime.DefaultHeaderPrefix = newPrefix
}

// GetDefaultHeaderPrefix set default header prefix
func GetDefaultHeaderPrefix() string {
	return mruntime.DefaultHeaderPrefix
}

// New new application
func New() application.Application {
	_initApp.Do(func() {
		_app = &Application{
			logger: golog.Default,
			config: conf.NewSetting(_configPath),
		}

		appPrefix := fmt.Sprintf("[%v@%v]", util.GetServiceName(_app.config.GetProjectName()), _app.config.GetProjectVersion())
		_app.logger.SetPrefix(appPrefix)
		_app.logger.SetTimeFormat("2006-01-02 15:04:05.000")

		_app.contextPool = application.New(func() application.Context {
			return application.NewContext(_app)
		})

		_app.controllerPool = application.NewControllerPool()
		_app.containerPool = application.NewContainerPool()
	})
	return _app
}

// DefaultGRPC default grpc server
func DefaultGRPC() application.Application {
	app := New()
	app.UseInterceptor(recovery.New(app))
	app.UseInterceptor(runtime.New(app))
	app.UseInterceptor(logger.New(app))
	app.UseInterceptor(di.New(app))
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
	return app.logger.Clone()
}

// RuntimeKeys get runtime keys
func (app *Application) RuntimeKeys() []truntime.RuntimeKey {
	app.mu.RLock()
	defer app.mu.RUnlock()
	newRuntimeKey := make([]truntime.RuntimeKey, len(app.runtimeKeys))
	copy(newRuntimeKey, app.runtimeKeys)
	return newRuntimeKey
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

// ControllerPool get all serviice pool
func (app *Application) ControllerPool() *application.ControllerPool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.controllerPool
}

// ContainerPool get all serviice pool
func (app *Application) ContainerPool() *application.ContainerPool {
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

// BindContainer bind container
func BindContainer(containerName reflect.Type, containerPool *sync.Pool) {
	newContainer := bootingContainer{
		containerName: containerName,
		containerPool: containerPool,
	}
	_bootingContainers = append(_bootingContainers, newContainer)
}

func (app *Application) initContainerPool() {

	app.Logger().Info(fmt.Sprintf("booting installing container start"))
	for _, container := range _bootingContainers {
		app.containerPool.NewContainer(container.containerName, container.containerPool)
		app.Logger().Info(fmt.Sprintf("booting installing container : %v ...installed", container.containerName))
	}
	line := fmt.Sprintf("booting installed %v container pool successfully", len(app.containerPool.GetContainerType()))
	app.Logger().Info(line)
}

// BindController bind service
func BindController(controllerName string, controllerPool *sync.Pool, requestMaps ...*application.RequestMap) {
	newController := bootingController{
		controllerName: controllerName,
		controllerPool: controllerPool,
		requestMaps:    requestMaps,
	}
	_bootingControllers = append(_bootingControllers, newController)
}

func (app *Application) initControllerPool() {
	app.Logger().Info(fmt.Sprintf("booting installing controller start"))
	for _, controller := range _bootingControllers {
		if len(controller.requestMaps) == 0 {
			app.controllerPool.NewController(controller.controllerName, controller.controllerPool)
			app.Logger().Info(fmt.Sprintf("booting installing controller : %v ...installed", controller.controllerName))
			continue
		}
		for _, request := range controller.requestMaps {
			newControllerName := fmt.Sprintf("%v@%v%v%v", request.Method, app.Conf().GetAppBaseURL(), controller.controllerName, request.SubPath)
			app.controllerPool.NewController(newControllerName, controller.controllerPool)
			app.controllerPool.NewControllerFunc(newControllerName, request.FuncName)
			app.controllerPool.NewControllerValidators(newControllerName, request.Validators...)
			realControllerName := strings.Replace(newControllerName, "@", " ==> ", -1)
			if app.controllerPool.ControllSelfCheck(newControllerName) {
				app.Logger().Info(fmt.Sprintf("booting installing controller : %v  -> %v ...installed", realControllerName, request.FuncName))
			} else {
				app.Logger().Info(fmt.Sprintf("booting installing controller : %v  -> no func registered, will use %v as default ...installed", realControllerName, request.Method))
			}
		}
	}
	line := fmt.Sprintf("booting installed %v controller pool successfully", len(app.controllerPool.GetControllerMap()))
	app.Logger().Info(line)
}

func (app *Application) initPool() {
	app.initControllerPool()
	app.initContainerPool()
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

func (app *Application) initHealthCheck() {
	app.router.GET(app.Conf().GetAppBaseURL()+_healthCheckPath, _defaultHealthCheckHandler(app))

	app.Logger().Info(fmt.Sprintf("booting installing healthChecker : %v  -> %v ...installed", app.Conf().GetAppBaseURL()+_healthCheckPath, "_defaultHealthCheckHandler"))
}

// InitRouter init router
// use gin framework by default
func (app *Application) InitRouter() {
	gin.DefaultWriter = ioutil.Discard
	app.router = gin.New()
	if app.Conf().GetCorsEnable() {
		app.router.Use(cors.New(cors.Config{
			AllowOrigins:     app.Conf().GetAllowOrigins(),
			AllowMethods:     app.Conf().GetAllowMethods(),
			AllowHeaders:     app.Conf().GetAllowHeaders(),
			ExposeHeaders:    app.Conf().GetExposeHeaders(),
			AllowCredentials: app.Conf().GetAllowCredentials(),
			MaxAge:           time.Duration(app.Conf().GetMaxAgeHour()) * time.Hour,
		}))
	}
	app.router.RedirectTrailingSlash = false
	if _enableHealthCheckPath {
		app.initHealthCheck()
	}
	app.router.GET(app.Conf().GetAppBaseURL()+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	for _, v := range app.middlewares {
		app.router.Use(v)
	}
	for _, controllerName := range app.ControllerPool().GetControllerMap() {
		controllerNameList := strings.Split(controllerName, "@")
		app.router.Handle(controllerNameList[0], controllerNameList[1], mdi.New(app))
	}
}

// Router get router
func (app *Application) Router() *gin.Engine {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.router
}

// InitHTTP serve grpc
func (app *Application) InitHTTP() {
	app.InitTrinity()
	app.InitRouter()
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
			line := fmt.Sprintf("booting http service registered successfully !")
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
	addr := fmt.Sprintf(":%v", app.config.GetAppPort())
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
