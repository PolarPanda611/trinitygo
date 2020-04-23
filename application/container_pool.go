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
	poolTags      map[string]reflect.Type
}

// NewContainerPool new pool with init map
func NewContainerPool() *ContainerPool {
	result := new(ContainerPool)
	result.poolMap = make(map[reflect.Type]*sync.Pool)
	result.poolTags = make(map[string]reflect.Type)
	return result

}

// NewContainer add new container
func (s *ContainerPool) NewContainer(containerType reflect.Type, containerPool *sync.Pool, containerTags []string) {
	s.poolMap[containerType] = containerPool
	s.containerType = append(s.containerType, containerType)
	if len(containerTags) > 0 {
		if containerTags[0] != "" {
			s.poolTags[containerTags[0]] = containerType
		}
	}
}

// GetContainerType get all service type
func (s *ContainerPool) GetContainerType(tags string) []reflect.Type {
	if tags != "" {
		var types []reflect.Type
		types = append(types, s.poolTags[tags])
		return types
	}
	return s.containerType
}

// GetContainer get service with di
func (s *ContainerPool) GetContainer(containerType reflect.Type, tctx Context, app Application, c *gin.Context) (interface{}, map[reflect.Type]interface{}) {
	pool, ok := s.poolMap[containerType]
	if !ok {
		panic("unknown service name")
	}
	service := pool.Get()
	sharedInstance := DiAllFields(service, tctx, app, c, false)
	return service, sharedInstance
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
