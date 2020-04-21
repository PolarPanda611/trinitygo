package application

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// DiAllFields di service pool
func DiAllFields(dest interface{}, tctx Context, app Application, c *gin.Context) map[reflect.Type]interface{} {
	sharedInstance := make(map[reflect.Type]interface{})
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if val.Kind() == reflect.Ptr && val.IsZero() {
			if !val.CanSet() {
				continue
			}
			if reflect.TypeOf(c) == val.Type() {
				val.Set(reflect.ValueOf(c))
				continue
			}
		}
		if val.Kind() == reflect.Interface && val.IsZero() {
			if !val.CanSet() {
				continue
			}
			if reflect.TypeOf(tctx).Implements(val.Type()) {
				isTransaction := reflect.TypeOf(dest).Elem().Field(index).Tag.Get("transaction")
				enableTx := false
				if app.Conf().GetAtomicRequest() {
					enableTx = true
				}
				if isTransaction != "" {
					if isTransaction == "true" {
						enableTx = true
					} else {
						enableTx = false
					}
				}
				if enableTx {
					tctx.DBTx()
				}
				val.Set(reflect.ValueOf(tctx))
				continue
			}
			for _, v := range app.ContainerPool().GetContainerType() {
				if v.Implements(val.Type()) {
					if instance, exist := sharedInstance[val.Type()]; exist {
						val.Set(reflect.ValueOf(instance))
						break
					}
					repo, sharedInstanceMap := app.ContainerPool().GetContainer(v, tctx, app, c)
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
