package service

import (
	"errors"
	"io/ioutil"
	"reflect"
	"sync"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/gin-gonic/gin"
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
	C    *gin.Context
	TCtx application.Context
}

func (s *ServiceImpl) Get(id int) (interface{}, error) {
	request, err := ioutil.ReadAll(s.C.Request.Body)
	if err != nil {
		return nil, err
	}
	if id < 100 {
		return string(request), nil
	}
	return nil, errors.New("wrong")
}
