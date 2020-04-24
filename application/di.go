package application

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
)

// DiSelfCheck ()map[reflect.Type]interface{} {}
func DiSelfCheck(destName interface{}, pool *sync.Pool, app Application) {
	controller := pool.Get()
	defer pool.Put(controller)
	controllerVal := reflect.Indirect(reflect.ValueOf(controller))
	app.Logger().Infof("booting self DI checking %v => %v ",
		destName,
		reflect.TypeOf(controller))
	for index := 0; index < controllerVal.NumField(); index++ {
		availableInjectInstance := 0
		var availableInjectType []reflect.Type
		val := controllerVal.Field(index)
		if !GetAutowiredTags(controller, index) {
			app.Logger().Infof("booting self DI checking Inject param: %v ,type:%v ,autowired tag not set, skipped ",
				reflect.TypeOf(controller).Elem().Field(index).Name,
				val.Type())
			continue
		}
		if !val.CanSet() {
			app.Logger().Fatalf("booting self DI checking Inject param: %v ,type: %v , private param , ...inject failed",
				reflect.TypeOf(controller).Elem().Field(index).Name,
				val.Type())
			continue
		}
		if !val.IsZero() {
			app.Logger().Fatalf("booting self DI checking Inject param: %v ,type:%v , not null param , ...inject failed",
				reflect.TypeOf(controller).Elem().Field(index).Name,
				val.Type())
			continue
		}
		if val.Kind() == reflect.Struct {
			app.Logger().Fatalf("booting self DI checking Inject param: %v ,type:%v , should be addressable , ...inject failed",
				reflect.TypeOf(controller).Elem().Field(index).Name,
				val.Type())
			continue
		}
		if val.Kind() == reflect.Ptr {
			// if is the gin context
			if reflect.TypeOf(&gin.Context{}) == val.Type() {
				app.Logger().Warn("The Gin Context already included in Trinity Go Context , You don't need to inject again !")
				app.Logger().Warnf("booting self DI checking Inject param: %v ,type:%v , injected ",
					reflect.TypeOf(controller).Elem().Field(index).Name,
					val.Type())
				continue
			}
			for _, v := range app.InstancePool().GetInstanceType(GetResourceTags(controller, index)) {
				if val.Type() == v {
					availableInjectType = append(availableInjectType, v)
					availableInjectInstance++
				}
			}
			availableInstanceLogger(availableInjectInstance, controller, index, val, app.Logger(), availableInjectType)
			continue
		}
		if val.Kind() == reflect.Interface {
			if reflect.TypeOf(&ContextImpl{}).Implements(val.Type()) {
				app.Logger().Infof("booting self DI checking Inject param: %v ,type:%v , injected ",
					reflect.TypeOf(controller).Elem().Field(index).Name,
					val.Type())
				continue
			}
			for _, v := range app.InstancePool().GetInstanceType(GetResourceTags(controller, index)) {
				if v.Implements(val.Type()) {
					availableInjectType = append(availableInjectType, v)
					availableInjectInstance++
				}
			}
			availableInstanceLogger(availableInjectInstance, controller, index, val, app.Logger(), availableInjectType)
			continue
		}
	}
}

// DiAllFields di service pool
func DiAllFields(dest interface{}, tctx Context, app Application, c *gin.Context) map[reflect.Type]interface{} {
	sharedInstance := make(map[reflect.Type]interface{})
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if !GetAutowiredTags(dest, index) {
			continue
		}
		if !val.CanSet() {
			continue
		}
		if !val.IsZero() {
			continue
		}
		if val.Kind() == reflect.Ptr {
			// if is the gin context
			if reflect.TypeOf(c) == val.Type() {
				val.Set(reflect.ValueOf(c))
				continue
			}
			for _, v := range app.InstancePool().GetInstanceType(GetResourceTags(dest, index)) {
				if val.Type() == v {
					if instance, exist := sharedInstance[val.Type()]; exist {
						val.Set(reflect.ValueOf(instance))
						break
					}
					repo, sharedInstanceMap := app.InstancePool().GetInstance(v, tctx, app, c)
					for instanceType, instanceValue := range sharedInstanceMap {
						sharedInstance[instanceType] = instanceValue
					}
					val.Set(reflect.ValueOf(repo))
					sharedInstance[val.Type()] = repo
					break
				}
			}
		}
		if val.Kind() == reflect.Interface {
			if reflect.TypeOf(tctx).Implements(val.Type()) {
				if !tctx.DBTxIsOpen() {
					enableTx := false
					if app.Conf().GetAtomicRequest() {
						enableTx = true
					}
					if TransactionTag(dest, index) {
						enableTx = true
					}
					if enableTx {
						tctx.DBTx()
					}
				}
				val.Set(reflect.ValueOf(tctx))
				continue
			}
			for _, v := range app.InstancePool().GetInstanceType(GetResourceTags(dest, index)) {
				if v.Implements(val.Type()) {
					if instance, exist := sharedInstance[val.Type()]; exist {
						val.Set(reflect.ValueOf(instance))
						break
					}
					repo, sharedInstanceMap := app.InstancePool().GetInstance(v, tctx, app, c)
					for instanceType, instanceValue := range sharedInstanceMap {
						sharedInstance[instanceType] = instanceValue
					}
					val.Set(reflect.ValueOf(repo))
					sharedInstance[val.Type()] = repo
					break
				}
			}
		}
	}
	return sharedInstance
}

// DiFree di controller
func DiFree(dest interface{}) {
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if !val.CanSet() {
			continue
		}
		if !val.IsZero() {
			val.Set(reflect.Zero(val.Type()))
		}
	}
}

// TransactionTag  get the transaction tag from struct
func TransactionTag(object interface{}, index int) bool {
	objectType := reflect.TypeOf(object)
	var isTransactionString string
	if objectType.Kind() == reflect.Struct {
		isTransactionString = reflect.TypeOf(object).Field(index).Tag.Get("autowired")
	} else {
		isTransactionString = reflect.TypeOf(object).Elem().Field(index).Tag.Get("autowired")
	}
	isTransaction, _ := strconv.ParseBool(isTransactionString)
	return isTransaction

}

// GetAutowiredTags get autowired tags from struct
func GetAutowiredTags(object interface{}, index int) bool {
	objectType := reflect.TypeOf(object)
	var isAutowiredString string
	if objectType.Kind() == reflect.Struct {
		isAutowiredString = objectType.Field(index).Tag.Get("autowired")
	} else {
		isAutowiredString = objectType.Elem().Field(index).Tag.Get("autowired")
	}
	isAutowired, _ := strconv.ParseBool(isAutowiredString)
	return isAutowired
}

// GetResourceTags get resource tags
func GetResourceTags(object interface{}, index int) string {
	objectType := reflect.TypeOf(object)
	if objectType.Kind() == reflect.Struct {
		return reflect.TypeOf(object).Field(index).Tag.Get("resource")
	}
	return reflect.TypeOf(object).Elem().Field(index).Tag.Get("resource")

}

func availableInstanceLogger(availableInjectInstance int, dest interface{}, index int, val reflect.Value, logger *golog.Logger, availableInjectType []reflect.Type) {
	if availableInjectInstance == 1 {
		logger.Infof("booting self DI checking Inject param: %v ,type:%v ,instance: %v ...injected ",
			reflect.TypeOf(dest).Elem().Field(index).Name,
			val.Type(),
			availableInjectType[0],
		)
	} else if availableInjectInstance < 1 {
		logger.Fatalf("booting self DI checking Inject param: %v ,type:%v , no instance available , ...inject failed",
			reflect.TypeOf(dest).Elem().Field(index).Name,
			val.Type())
	} else {
		availableType := ""
		for _, v := range availableInjectType {
			availableType += fmt.Sprintf("%v,", v)
		}
		logger.Fatalf("booting self DI checking Inject param: %v ,type:%v , more than one instance (%v) available , ...inject failed",
			reflect.TypeOf(dest).Elem().Field(index).Name,
			val.Type(),
			availableType)
	}
}
