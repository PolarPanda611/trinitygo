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
	if _, ok := s.poolMap[containerType]; ok {
		return
	}
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
		if container, ok := s.poolTags[tags]; ok {
			types = append(types, container)
			return types
		}
		return types
	}
	return s.containerType
}

// CheckContainerNameIfExist check contain name if exist
func (s *ContainerPool) CheckContainerNameIfExist(containerName reflect.Type) bool {
	_, ok := s.poolMap[containerName]
	return ok
}

// ContainerDISelfCheck  self check di request registered func exist or not
func (s *ContainerPool) ContainerDISelfCheck(app Application) {
	for controllerName, pool := range s.poolMap {
		app.Logger().Infof("booting self DI checking container %v ", controllerName)

		DiSelfCheck(controllerName, pool, app)
	}
	return

}

// GetContainer get service with di
func (s *ContainerPool) GetContainer(containerType reflect.Type, tctx Context, app Application, c *gin.Context) (interface{}, map[reflect.Type]interface{}) {
	pool, ok := s.poolMap[containerType]
	if !ok {
		panic("unknown service name")
	}
	service := pool.Get()
	sharedInstance := DiAllFields(service, tctx, app, c)
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
