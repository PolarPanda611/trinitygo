package service

import (
	"errors"
	"reflect"
	"sync"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
)

func init() {
	trinitygo.BindContainer(reflect.TypeOf(&ServiceImpl{}), &sync.Pool{
		New: func() interface{} {
			service := new(ServiceImpl)
			return service
		},
	})
}

type Service interface {
	Get(id int) (interface{}, error)
}
type ServiceImpl struct {
	TCtx application.Context
}

func (s *ServiceImpl) Get(id int) (interface{}, error) {
	if id < 100 {
		return s.TCtx.GetRuntime(), nil
	}
	return nil, errors.New("wrong")
}
