package application

import (
	"fmt"
	"reflect"
)

// DiServicePool di service pool
func DiServicePool(dest interface{}, tctx Context, app Application) []interface{} {
	var toFreeRepository []interface{}
	DiTCtx(dest, tctx)
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if val.Kind() == reflect.Interface && val.IsZero() {

			for _, v := range app.GetRepositoryPool().GetRepositoryType() {
				if !v.Implements(val.Type()) {
					continue
				}
				if !val.CanSet() {
					// not the public param , cannot inject
					continue
				}
				repo := app.GetRepositoryPool().GetRepository(v, tctx, app)
				toFreeRepository = append(toFreeRepository, repo)
				val.Set(reflect.ValueOf(repo))
				continue
			}
		}
	}
	return toFreeRepository
}

// DiController di controller
func DiController(dest interface{}, tCtx Context, app Application) ([]interface{}, []interface{}) {
	var toFreeService []interface{}
	var toFreeRepository []interface{}
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if val.Kind() == reflect.Interface && val.IsZero() {
			for _, v := range app.GetServicePool().GetServiceType() {
				if !v.Implements(val.Type()) {
					continue
				}
				if !val.CanSet() {
					fmt.Println("must be public")
					continue
				}
				service, toFreeRepo := app.GetServicePool().GetService(v, tCtx, app)
				toFreeService = append(toFreeService, service)
				toFreeRepository = append(toFreeRepository, toFreeRepo...)
				val.Set(reflect.ValueOf(service))

			}
		}
	}
	return toFreeService, toFreeRepository
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

// DiTCtx di trintiy context
func DiTCtx(dest interface{}, tCtx Context) {
	destVal := reflect.Indirect(reflect.ValueOf(dest))
	for index := 0; index < destVal.NumField(); index++ {
		val := destVal.Field(index)
		if val.Kind() == reflect.Interface && val.IsZero() {

			if !reflect.TypeOf(tCtx).Implements(val.Type()) {
				// not implement
				continue
			}
			if !val.CanSet() {
				// not the public param , cannot inject
				continue
			}
			val.Set(reflect.ValueOf(tCtx))
		}
	}
}
