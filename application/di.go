package application

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

// DiAllFields di service pool
func DiAllFields(dest interface{}, tctx Context, app Application, c *gin.Context) []interface{} {
	var toFreeContainer []interface{}
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if val.Kind() == reflect.Interface && val.IsZero() {
			if !val.CanSet() {
				// not the public param , cannot inject
				continue
			}

			// check if implement tctx
			if reflect.TypeOf(tctx).Implements(val.Type()) {
				// if  implemented
				val.Set(reflect.ValueOf(tctx))
				continue
			}

			// check if implement gin.context
			if reflect.TypeOf(c).Implements(val.Type()) {
				// if  implemented
				val.Set(reflect.ValueOf(c))
				continue
			}

			for _, v := range app.GetContainerPool().GetContainerType() {
				if v.Implements(val.Type()) {
					repo, subToFreeContainer := app.GetContainerPool().GetContainer(v, tctx, app, c)
					toFreeContainer = append(toFreeContainer, repo)
					toFreeContainer = append(toFreeContainer, subToFreeContainer...)
					val.Set(reflect.ValueOf(repo))
					continue
				}
			}
		}
	}
	return toFreeContainer
}

// DiFree di controller
func DiFree(dest interface{}) {
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if val.Kind() == reflect.Interface {
			if !val.IsZero() {
				val.Set(reflect.Zero(val.Type()))
			}
		}
	}
}
