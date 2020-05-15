package application

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
)

type diSelfCheckResultCount struct {
	info    int
	warning int
}

// DiSelfCheck ()map[reflect.Type]interface{} {}
func DiSelfCheck(destName interface{}, pool *sync.Pool, logger *golog.Logger, instancePool *InstancePool, instanceMapping map[string]reflect.Type) {
	resultCount := new(diSelfCheckResultCount)
	instance := pool.Get()
	defer pool.Put(instance)
	instanceVal := reflect.Indirect(reflect.ValueOf(instance))
	for index := 0; index < instanceVal.NumField(); index++ {
		objectName := encodeObjectName(instance, index)
		availableInjectInstance := 0
		var availableInjectType []reflect.Type
		val := instanceVal.Field(index)
		if !GetAutowiredTags(instance, index) {
			if val.CanSet() {
				logDIWarnf(objectName, logger, resultCount, "autowired tag not set")
			}
			continue
		}
		if !val.CanSet() {
			logDIFatalf(objectName, logger, resultCount, "private param")
			continue
		}
		if !val.IsZero() {
			logDIFatalf(objectName, logger, resultCount, "not null param")
			continue
		}
		if val.Kind() == reflect.Struct {
			logDIFatalf(objectName, logger, resultCount, "should be addressable")
			continue
		}
		if val.Kind() == reflect.Ptr {
			// if is the gin context
			if reflect.TypeOf(&gin.Context{}) == val.Type() {
				logDIWarnf(objectName, logger, resultCount, "Should use Application.Context instead of gin.Context !")
				logDIInfof(objectName, logger, resultCount, reflect.TypeOf(&gin.Context{}), instanceMapping)
				continue
			}
			for _, v := range instancePool.GetInstanceType(GetResourceTags(instance, index)) {
				if val.Type() == v {
					availableInjectType = append(availableInjectType, v)
					availableInjectInstance++
				}
			}
			availableInstanceLogger(availableInjectInstance, objectName, logger, resultCount, availableInjectType, instanceMapping)
			continue
		}
		if val.Kind() == reflect.Interface {
			if reflect.TypeOf(&ContextImpl{}).Implements(val.Type()) {
				logDIInfof(objectName, logger, resultCount, reflect.TypeOf(&ContextImpl{}), instanceMapping)
				continue
			}
			for _, v := range instancePool.GetInstanceType(GetResourceTags(instance, index)) {
				if v.Implements(val.Type()) {
					availableInjectType = append(availableInjectType, v)
					availableInjectInstance++
				}
			}
			availableInstanceLogger(availableInjectInstance, objectName, logger, resultCount, availableInjectType, instanceMapping)
			continue
		}
	}
	logger.Infof("booting self DI checking instance %v finished with info %v , warning %v", destName, resultCount.info, resultCount.warning)
}

// DiAllFields di service pool
func DiAllFields(dest interface{}, tctx Context, app Application, c *gin.Context, instanceMapping map[string]reflect.Type) map[reflect.Type]interface{} {
	sharedInstance := make(map[reflect.Type]interface{})
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if !GetAutowiredTags(dest, index) {
			continue
		}
		objectName := encodeObjectName(dest, index)
		switch instanceMapping[objectName] {
		case reflect.TypeOf(c):
			val.Set(reflect.ValueOf(c))
			break
		case reflect.TypeOf(tctx):
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
			break
		default:
			if instanceMapping[objectName] == reflect.TypeOf(c) {
				val.Set(reflect.ValueOf(c))
				continue
			}

			repo, sharedInstanceMap := app.InstancePool().GetInstance(instanceMapping[objectName], tctx, app, c)
			for instanceType, instanceValue := range sharedInstanceMap {
				sharedInstance[instanceType] = instanceValue
			}
			val.Set(reflect.ValueOf(repo))
			sharedInstance[val.Type()] = repo
			break
		}
	}
	return sharedInstance
}

// DiFree di instance
func DiFree(dest interface{}) {
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

func availableInstanceLogger(availableInjectInstanceCount int, objectName string, logger *golog.Logger, resultCount *diSelfCheckResultCount, availableInjectType []reflect.Type, instanceMapping map[string]reflect.Type) {
	if availableInjectInstanceCount == 1 {
		logDIInfof(objectName, logger, resultCount, availableInjectType[0], instanceMapping)
	} else if availableInjectInstanceCount < 1 {
		logDIFatalf(objectName, logger, resultCount, "no instance available")
	} else {
		availableType := ""
		for _, v := range availableInjectType {
			availableType += fmt.Sprintf("%v,", v)
		}
		logDIFatalf(objectName, logger, resultCount, fmt.Sprintf("more than one instance (%v) available", availableType))
	}
}

func encodeObjectName(instance interface{}, index int) string {
	instanceType := reflect.TypeOf(instance)
	instanceVal := reflect.Indirect(reflect.ValueOf(instance))
	paramName := instanceType.Elem().Field(index).Name
	paramType := instanceVal.Field(index).Type()
	return fmt.Sprintf("%v.%v.(%v)", instanceType, paramName, paramType)

}
func logDIInfof(objectName string, logger *golog.Logger, resultCount *diSelfCheckResultCount, instance reflect.Type, instanceMapping map[string]reflect.Type) {
	logger.Infof("booting self DI checking Inject object: %v, with instance: %v,  ...injected ", objectName, instance)
	resultCount.info++
	instanceMapping[objectName] = instance
	return
}

func logDIWarnf(objectName string, logger *golog.Logger, resultCount *diSelfCheckResultCount, errString string) {
	logger.Warnf("booting self DI checking Inject object: %v, %v, skipped ", objectName, errString)
	resultCount.warning++
}

func logDIFatalf(objectName string, logger *golog.Logger, resultCount *diSelfCheckResultCount, errString string) {
	logger.Fatalf("booting self DI checking Inject object: %v, %v, ...inject failed ", objectName, errString)
}
