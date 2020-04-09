package application

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
)

// ContainerPool service pool
type ContainerPool struct {
	poolMap       map[reflect.Type]*sync.Pool
	containerType []reflect.Type
}

// NewContainerPool new pool with init map
func NewContainerPool() *ContainerPool {
	result := new(ContainerPool)
	result.poolMap = make(map[reflect.Type]*sync.Pool)
	return result

}

// NewContainer add new container
func (s *ContainerPool) NewContainer(containerType reflect.Type, containerPool *sync.Pool) {
	s.poolMap[containerType] = containerPool
	s.containerType = append(s.containerType, containerType)
}

// GetContainerType get all service type
func (s *ContainerPool) GetContainerType() []reflect.Type {
	return s.containerType
}

// GetContainer get service with di
func (s *ContainerPool) GetContainer(containerType reflect.Type, tctx Context, app Application, c *gin.Context) (interface{}, []interface{}) {
	pool, ok := s.poolMap[containerType]
	if !ok {
		panic("unknown service name")
	}
	service := pool.Get()
	toFreeContainer := DiContainerPool(service, tctx, app, c)
	return service, toFreeContainer
}

// Release release service
func (s *ContainerPool) Release(container interface{}) {
	t := reflect.TypeOf(container)
	syncpool, ok := s.poolMap[t]
	if !ok {
		return
	}
	DiFree(container)
	syncpool.Put(container)
}
