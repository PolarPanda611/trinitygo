package trinitygo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/PolarPanda611/bar"
	"github.com/PolarPanda611/trinitygo/httputil"
	"github.com/PolarPanda611/trinitygo/keyword"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/PolarPanda611/trinitygo/startup"
	"github.com/PolarPanda611/trinitygo/util"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
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
	_startupLatency        int64 = 0
	_responseFactory       func(status int, res interface{}, runtime map[string]string) interface{}
	_                      application.Application = new(Application)
	_casbinEnable          bool                    = false
	_app                   *Application
	_initApp               sync.Once
	_configPath            string = "./config/"
	_casbinConfPath        string = "./config/casbin.conf"
	_bootingControllers    []bootingController
	_bootingInstances      []bootingInstance
	_bootingModels         []bootingModel
	_healthCheckPath       string                                                                              = "/ping"
	_enableHealthCheckPath bool                                                                                = false
	_logSelfCheck          bool                                                                                = true
	_funcToGetWhoAmI       func(app application.Application, c *gin.Context, db *gorm.DB) (interface{}, error) = func(app application.Application, c *gin.Context, db *gorm.DB) (interface{}, error) {
		return nil, nil
	}
	_defaultHealthCheckHandler func(app application.Application) func(c *gin.Context) = func(app application.Application) func(c *gin.Context) {
		return func(c *gin.Context) {
			err := app.DB().DB().Ping()
			if err != nil {
				c.AbortWithStatusJSON(400, httputil.ResponseData{
					Status: 400,
					Err:    err,
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

type bootingModel struct {
	modelInstance interface{}
	defaultValues []interface{}
}

type bootingController struct {
	controllerName string
	instance       interface{}
	requestMaps    []*application.RequestMap
}

type bootingInstance struct {
	instanceName reflect.Type
	instancePool *sync.Pool
	instanceTags []string
}

// Application core of trinity
type Application struct {
	b            *bar.Bar
	config       conf.Conf
	keyword      keyword.Keyword
	logger       *golog.Logger
	logSelfCheck bool
	contextPool  *application.ContextPool

	// used for build
	once           sync.Once
	mu             sync.RWMutex
	db             *gorm.DB
	enforcer       *casbin.Enforcer
	controllerPool *application.ControllerPool
	instancePool   *application.InstancePool

	//grpc
	interceptors []grpc.UnaryServerInterceptor
	serviceMesh  sd.ServiceMesh
	grpcServer   *grpc.Server

	//http
	middlewares []gin.HandlerFunc
	router      *gin.Engine

	//runtime
	runtimeKeys     []truntime.RuntimeKey
	responseFactory func(status int, res interface{}, runtime map[string]string) interface{}
}

// SetKeyword set keywork list
func SetKeyword(k keyword.Keyword) {
	keyword.SetKeyword(k)
}

// SetConfigPath set config path
func SetConfigPath(path string) {
	_configPath = path
}

// SetResponseFactory set response factory
func SetResponseFactory(f func(status int, res interface{}, runtime map[string]string) interface{}) {
	_responseFactory = f
}

// SetCasbinConfPath set config path
func SetCasbinConfPath(path string) {
	_casbinConfPath = path
}

// SetFuncGetWhoAmI set _funcToGetWhoAmI
func SetFuncGetWhoAmI(f func(app application.Application, c *gin.Context, db *gorm.DB) (interface{}, error)) {
	_funcToGetWhoAmI = f
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

// SetIsLogSelfCheck set is log self check
// by default , it is true
func SetIsLogSelfCheck(isLog bool) {
	_logSelfCheck = isLog
}

// New new application
func New(args ...bool) application.Application {
	_initApp.Do(func() {
		_app = &Application{
			logger:       golog.Default,
			config:       conf.NewSetting(_configPath),
			logSelfCheck: _logSelfCheck,
			keyword:      keyword.GetKeyword(),
		}
		_app.b = bar.NewBar(0, 100)
		appPrefix := fmt.Sprintf("[%v@%v]", util.GetServiceName(_app.config.GetProjectName()), _app.config.GetProjectVersion())
		_app.logger.SetPrefix(appPrefix)
		_app.logger.SetTimeFormat("2006-01-02 15:04:05.000")
		_app.setProgress(2, _startupLatency, "init logger")
		_app.contextPool = application.New(func() application.Context {
			return application.NewContext(_app, _funcToGetWhoAmI)
		})
		_app.setProgress(4, _startupLatency, "init context pool ")
		_app.controllerPool = application.NewControllerPool()
		_app.setProgress(6, _startupLatency, "init controller pool ")
		_app.instancePool = application.NewInstancePool()
		_app.setProgress(8, _startupLatency, "init instance pool ")
		if _responseFactory != nil {
			_app.responseFactory = _responseFactory
		}
		_app.setProgress(10, _startupLatency, "init response factory ")
		if len(args) != 0 {
			startup.SetStartupDebugger(args[0])
		}
	})
	return _app
}

// DefaultGRPC default grpc server
func DefaultGRPC() application.Application {
	app := New()
	app.UseInterceptor(logger.New(app))
	app.UseInterceptor(recovery.New(app))
	app.UseInterceptor(runtime.New(app))
	app.UseInterceptor(di.New(app))
	return app
}

// DefaultHTTP default http server
func DefaultHTTP(args ...bool) application.Application {
	app := New(args...)

	return app
}

// Logger get logger
func (app *Application) setProgress(progress, latency int64, message string) {
	time.Sleep(time.Duration(latency) * time.Millisecond)
	app.b.Cur <- bar.CurrentStep{
		Cur:      progress,
		Message:  message,
		IsReturn: false,
	}
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

// Keyword get key word list
func (app *Application) Keyword() keyword.Keyword {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.keyword
}

// DB get db
func (app *Application) DB() *gorm.DB {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.db
}

// ResponseFactory get response factory
func (app *Application) ResponseFactory() func(status int, res interface{}, runtime map[string]string) interface{} {
	app.mu.RLock()
	defer app.mu.RUnlock()
	if app.responseFactory != nil {
		return app.responseFactory
	}
	return nil
}

// ContextPool get contextpoo;
func (app *Application) ContextPool() *application.ContextPool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.contextPool
}

// ControllerPool get all service pool
func (app *Application) ControllerPool() *application.ControllerPool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.controllerPool
}

// InstancePool get all service pool
func (app *Application) InstancePool() *application.InstancePool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.instancePool
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
	startup.AppendStartupDebuggerInfo("booting detected db initializer ...install successfully , test passed ! ")
}

// IsLogSelfCheck if log self check
func (app *Application) IsLogSelfCheck() bool {
	return app.logSelfCheck
}

// RegisterModel register model
// @instance should be ptr instance
// @defaultValues should be ptr instance
// default value will create the instance and will not update .
func RegisterModel(instance interface{}, defaultValues ...interface{}) {
	newModel := bootingModel{
		modelInstance: instance,
	}
	newModel.defaultValues = defaultValues
	_bootingModels = append(_bootingModels, newModel)
}

// RegisterController register service
// @controllerName basic request path of the controller
// @instance controller instance  ,should be ptr
// @requestMaps optional ,register the request based on the basic path of controller
func RegisterController(controllerName string, instance interface{}, requestMaps ...*application.RequestMap) {
	newController := bootingController{
		controllerName: controllerName,
		instance:       instance,
		requestMaps:    requestMaps,
	}
	_bootingControllers = append(_bootingControllers, newController)
}

// RegisterInstance register instance
// @instance instance of the class , should be ptr
// @tags the tags for the instance to specify
func RegisterInstance(instance interface{}, tags ...string) {
	var newInstance bootingInstance
	switch reflect.TypeOf(instance).Kind() {
	case reflect.Struct:
		newInstance = bootingInstance{
			instanceName: reflect.New(reflect.TypeOf(instance)).Type(),
			instancePool: &sync.Pool{
				New: func() interface{} {
					return reflect.New(reflect.TypeOf(instance)).Interface()
				},
			},
			instanceTags: tags,
		}
		break
	case reflect.Func:
		f, ok := instance.(func() interface{})
		if !ok {
			log.Fatal("The instance func should be  func () interface{}")
		}
		newInstance = bootingInstance{
			instanceName: reflect.TypeOf(f()),
			instancePool: &sync.Pool{
				New: func() interface{} {
					return f()
				},
			},
			instanceTags: tags,
		}
		break
	default:
		log.Fatal("The instance should be struct or func () interface{}")
	}
	_bootingInstances = append(_bootingInstances, newInstance)
}

// init model and init default value
func (app *Application) initModel() {
	startup.AppendStartupDebuggerInfo("booting installing model start")
	// no model register , skip
	if len(_bootingModels) == 0 {
		return
	}
	// db not initialize , skip
	if app.db == nil {
		return
	}
	// sync model structure to db
	for _, bootingModel := range _bootingModels {
		instanceKind := reflect.TypeOf(bootingModel.modelInstance).Kind()
		if instanceKind != reflect.Ptr {
			app.logger.Fatalf("booting installing model : %v failed , kind %v ,  should be ptr , ", reflect.TypeOf(bootingModel.modelInstance), instanceKind)
		}
		if err := app.db.AutoMigrate(bootingModel.modelInstance).Error; err != nil {
			app.logger.Fatalf("booting installing model : %v failed , err : %v , ", reflect.TypeOf(bootingModel.modelInstance), err)
		}
		startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting installing model : %v  , ...installed , ", reflect.TypeOf(bootingModel.modelInstance)))
		if bootingModel.defaultValues != nil {
			for _, defaultValue := range bootingModel.defaultValues {
				defaultValueKind := reflect.TypeOf(defaultValue).Kind()
				switch defaultValueKind {
				case reflect.Ptr:
					if err := app.db.FirstOrCreate(defaultValue, defaultValue).Error; err != nil {
						app.logger.Fatalf("booting installing model : %v , set default value %v  , err : %v , ", reflect.TypeOf(bootingModel.modelInstance), defaultValue, err)
					}
					break
				default:
					app.logger.Fatalf("booting installing model : %v , set default value %v  type %v , err : %v , ", reflect.TypeOf(bootingModel.modelInstance), defaultValue, defaultValueKind, "default value only support ptr ")
				}
			}
		}
	}
	startup.AppendStartupDebuggerInfo("booting installing model end")
	app.setProgress(30, _startupLatency, "init model")
}

// initControllerPool initial controller pool
func (app *Application) initControllerPool() {
	startup.AppendStartupDebuggerInfo("booting installing controller start")
	for _, controller := range _bootingControllers {
		if len(controller.requestMaps) == 0 {
			app.controllerPool.NewController(controller.controllerName, reflect.New(reflect.TypeOf(controller.instance)).Type())
			RegisterInstance(controller.instance, controller.controllerName)
			startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting installing controller : %v ...installed", controller.controllerName))
			continue
		}
		for _, request := range controller.requestMaps {
			newControllerName := fmt.Sprintf("%v@%v%v%v", request.Method, app.Conf().GetAppBaseURL(), controller.controllerName, request.SubPath)
			app.controllerPool.NewController(newControllerName, reflect.New(reflect.TypeOf(controller.instance)).Type())
			RegisterInstance(controller.instance, newControllerName)
			app.controllerPool.NewControllerFunc(newControllerName, request.FuncName)
			app.controllerPool.NewControllerValidators(newControllerName, request.Validators...)
			realControllerName := strings.Replace(newControllerName, "@", " ==> ", -1)
			startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting installing controller : %v  -> %v ...installed", realControllerName, request.FuncName))
			startup.AppendRequestMapping(strings.Split(newControllerName, "@")[0], strings.Split(newControllerName, "@")[1], request.FuncName)
		}
	}
	startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting installed %v controller pool successfully", len(app.controllerPool.GetControllerMap())))
	app.setProgress(40, _startupLatency, "init controller pool")
}

// initInstancePool initial instance pool
func (app *Application) initInstancePool() {
	for _, instance := range _bootingInstances {
		if len(instance.instanceTags) > 0 && len(app.instancePool.GetInstanceType(instance.instanceTags[0])) > 0 {
			app.Logger().Fatalf("booting installing instance : %v failed ,instance with tag %v already exists", instance.instanceName, instance.instanceTags[0])
		}
		if app.instancePool.CheckInstanceNameIfExist(instance.instanceName) {
			continue
		}
		app.instancePool.NewInstance(instance.instanceName, instance.instancePool, instance.instanceTags)
		startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting installing instance : %v ...installed", instance.instanceName))
	}
	startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting installed %v instance pool successfully", len(app.instancePool.GetInstanceType(""))))
	app.setProgress(50, _startupLatency, "init instance pool")
}

func (app *Application) initSelfCheck() {
	app.controllerPool.ControllerFuncSelfCheck(app.instancePool, app.IsLogSelfCheck(), app.logger)
	startup.AppendStartupDebuggerInfo("booting self func checking controller successfully ")
	app.instancePool.InstanceDISelfCheck(app)
	app.setProgress(60, _startupLatency, "init instance pool")
}
func (app *Application) initPool() {
	app.initModel()
	app.initControllerPool()
	app.initInstancePool()
	app.initSelfCheck()

}

func (app *Application) initRuntime() {
	runtimeKey := ""
	for _, v := range app.runtimeKeys {
		runtimeKey += fmt.Sprintf("%v,", v.GetKeyName())
	}
	// runtime checking
	startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting detected %v runtime([%v]) ...installed", len(app.runtimeKeys), runtimeKey))
	app.setProgress(12, _startupLatency, "init runtime")
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
		startup.AppendStartupDebuggerInfo("booting db instance with default install")
	}
	app.setProgress(14, _startupLatency, "init db")
}

func (app *Application) initCasbin() {
	if !_casbinEnable {
		return
	}
	adapter, err := gormadapter.NewAdapterByDBUsePrefix(app.db, app.Conf().GetDBTablePrefix())
	if err != nil {
		app.logger.Fatal("create casbin adapter err ", err)
	}
	app.enforcer, err = casbin.NewEnforcer(_casbinConfPath, adapter)
	if err != nil {
		app.logger.Fatal("create casbin enforcer err ", err)
	}
	err = app.enforcer.LoadPolicy()
	if err != nil {
		app.logger.Fatal("load casbin enforcer err ", err)
	}
	app.setProgress(16, _startupLatency, "init casbin")
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
				app.logger.Fatal("get service mesh client err", err)
			}
			app.serviceMesh = c
			break
		default:
			app.logger.Fatal("wrong service mash type")
		}
	}
	app.setProgress(65, _startupLatency, "init service discovery")
}

// InitTrinity serve grpc
func (app *Application) InitTrinity() {
	app.initRuntime()
	app.initDB()
	app.initCasbin()
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
	startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting detected %v interceptors ...installed ", len(app.interceptors)))
	app.grpcServer = grpc.NewServer(opts...)
}

func (app *Application) initHealthCheck() {
	app.router.GET(app.Conf().GetAppBaseURL()+_healthCheckPath, _defaultHealthCheckHandler(app))
	startup.AppendStartupDebuggerInfo(fmt.Sprintf("booting installing healthChecker : %v  -> %v ...installed", app.Conf().GetAppBaseURL()+_healthCheckPath, "_defaultHealthCheckHandler"))
	startup.AppendRequestMapping("GET", app.Conf().GetAppBaseURL()+_healthCheckPath, "_defaultHealthCheckHandler")
}

// Enforcer get casbin enforcer instance
func (app *Application) Enforcer() *casbin.Enforcer {
	if app.enforcer == nil {
		app.logger.Fatal("You need init casbin first ")
	}
	return app.enforcer
}

// InitRouter init router
// use gin framework by default
func (app *Application) InitRouter() {
	gin.DefaultWriter = ioutil.Discard
	app.router = gin.New()
	app.setProgress(70, _startupLatency, "init gin router")
	app.router.Use(mlogger.New(app))
	app.setProgress(75, _startupLatency, "init gin logger middleware")
	app.router.Use(httprecovery.New(app))
	app.setProgress(80, _startupLatency, "init gin recovery middleware")
	if app.Conf().GetCorsEnable() {
		app.router.Use(cors.New(cors.Config{
			AllowOrigins:     app.Conf().GetAllowOrigins(),
			AllowMethods:     app.Conf().GetAllowMethods(),
			AllowHeaders:     app.Conf().GetAllowHeaders(),
			ExposeHeaders:    app.Conf().GetExposeHeaders(),
			AllowCredentials: app.Conf().GetAllowCredentials(),
			MaxAge:           time.Duration(app.Conf().GetMaxAgeHour()) * time.Hour,
			AllowOriginFunc: func(origin string) bool {
				return strings.Contains(origin, "127.0.0.1") || strings.Contains(origin, "localhost") || strings.Contains(origin, "http://192")
			},
		}))
	}
	app.setProgress(83, _startupLatency, "init gin cors middleware")

	app.router.GET(app.Conf().GetAppBaseURL()+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	startup.AppendRequestMapping("GET", app.Conf().GetAppBaseURL()+"/swagger/*any")
	app.setProgress(85, _startupLatency, "init swagger docs handler")
	app.router.Static(app.Conf().GetAppBaseURL()+app.Conf().GetAppStaticURL(), app.Conf().GetAppStaticPath())
	startup.AppendRequestMapping("GET", app.Conf().GetAppBaseURL()+app.Conf().GetAppStaticURL())
	app.setProgress(87, _startupLatency, "init static file handler")
	app.router.Static(app.Conf().GetAppBaseURL()+app.Conf().GetAppMediaURL(), app.Conf().GetAppMediaPath())
	app.setProgress(89, _startupLatency, "init media file handler")
	startup.AppendRequestMapping("GET", app.Conf().GetAppBaseURL()+app.Conf().GetAppMediaURL())
	for _, v := range app.middlewares {
		app.router.Use(v)
	}
	app.router.Use(mruntime.New(app))
	app.setProgress(91, _startupLatency, "init runtime middleware")
	if _enableHealthCheckPath {
		app.initHealthCheck()
	}
	for _, v := range app.middlewares {
		app.router.Use(v)
	}
	for _, controllerName := range app.ControllerPool().GetControllerMap() {
		controllerNameList := strings.Split(controllerName, "@")
		app.router.Handle(controllerNameList[0], controllerNameList[1], mdi.New(app))
	}
	app.setProgress(95, _startupLatency, "init http handler")
}

// Router get router
func (app *Application) Router() *gin.Engine {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.router
}

// InitHTTP serve http
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
				app.config.GetServiceDiscoveryTimeout(),
			); err != nil {
				gErr <- err
			}
			startup.AppendStartupDebuggerInfo("booting http service registered successfully !")
			app.setProgress(97, _startupLatency, "register service in SD")
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
		app.setProgress(100, _startupLatency, fmt.Sprintf("http service listen at %v ", addr))
		startup.AppendStartupDebuggerInfo("\n" + util.GenerateFiglet(app.config.GetProjectName()))
		startup.AppendStartupDebuggerInfo(fmt.Sprintf("booted http service listen at %v started ", addr))
		fmt.Println()
		if startup.GetStartupDebugger() {

			for _, v := range startup.GetStartupDebuggerInfo() {
				app.Logger().Info(v)
			}
		}
		for _, v := range startup.GetRequestMapping() {
			app.Logger().Info(v)
		}
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
			app.config.GetServiceDiscoveryTimeout(),
		)
		if err != nil {
			line = fmt.Sprintf("booting http service deRegistered failed !")
			app.Logger().Error(line)
		} else {
			line = fmt.Sprintf("booting http service deRegistered successfully !")
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
				app.config.GetServiceDiscoveryTimeout(),
			); err != nil {
				gErr <- err
			}
			startup.AppendStartupDebuggerInfo("booting grpc service registered successfully !")
		}
		startup.AppendStartupDebuggerInfo("\n" + util.GenerateFiglet(app.config.GetProjectName()))
		startup.AppendStartupDebuggerInfo(fmt.Sprintf("booted grpc service listen at %v started", addr))
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
			app.config.GetServiceDiscoveryTimeout(),
		)
		if err != nil {
			line = fmt.Sprintf("booting grpc service deRegistered failed !")
			app.Logger().Error(line)
		} else {
			line = fmt.Sprintf("booting grpc service deRegistered successfully !")
			app.Logger().Info(line)
		}

	}

}
