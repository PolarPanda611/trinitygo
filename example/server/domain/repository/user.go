package repository

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/PolarPanda611/trinitygo"

	"github.com/PolarPanda611/trinitygo/application"
)

func init() {
	trinitygo.BindRepository(reflect.TypeOf(&UserRepoImpl{}), &sync.Pool{
		New: func() interface{} {
			service := new(UserRepoImpl)
			return service
		},
	})
}

// UserRepo user repo
type UserRepo interface {
	Print() string
}

// UserRepoImpl user repo impl
type UserRepoImpl struct {
	TContext application.Context
}

// Print user repo impl print method
func (r *UserRepoImpl) Print() string {
	return "userrepo" + fmt.Sprintf("%v", r.TContext.GetRuntime())
}
