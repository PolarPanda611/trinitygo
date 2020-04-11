package application

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
)

// DiAllFields di service pool
func DiAllFields(dest interface{}, tctx Context, app Application, c *gin.Context) []interface{} {
	var toFreeContainer []interface{}
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if val.Kind() == reflect.Ptr && val.IsZero() {
			if !val.CanSet() {
				fmt.Println("skip")
			}

			if reflect.TypeOf(c) == val.Type() {
				val.Set(reflect.ValueOf(c))
			}

		}
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

			for _, v := range app.ContainerPool().GetContainerType() {
				if v.Implements(val.Type()) {
					repo, subToFreeContainer := app.ContainerPool().GetContainer(v, tctx, app, c)
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
		if !val.IsZero() {
			// fmt.Println(val.Type(), "set null")
			val.Set(reflect.Zero(val.Type()))
		}
	}
}
