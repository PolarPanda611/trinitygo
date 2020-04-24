package application

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
)

// InstancePool service pool
type InstancePool struct {
	// map[instanceName] = instancePool
	poolMap map[reflect.Type]*sync.Pool
	// instance list
	instanceTypeList []reflect.Type
	// instance tags maps instanceName
	poolTags map[string]reflect.Type
	// instanceMapping instance mapping
	// caching the instance di instance relation during the self check
	instanceMapping map[reflect.Type]reflect.Type
}

// NewInstancePool new pool with init map
func NewInstancePool() *InstancePool {
	result := new(InstancePool)
	result.poolMap = make(map[reflect.Type]*sync.Pool)
	result.poolTags = make(map[string]reflect.Type)
	result.instanceMapping = make(map[reflect.Type]reflect.Type)
	return result

}

// NewInstance add new instance
func (s *InstancePool) NewInstance(instanceType reflect.Type, instancePool *sync.Pool, instanceTags []string) {
	if _, ok := s.poolMap[instanceType]; ok {
		return
	}
	s.poolMap[instanceType] = instancePool
	s.instanceTypeList = append(s.instanceTypeList, instanceType)
	if len(instanceTags) > 0 {
		if instanceTags[0] != "" {
			s.poolTags[instanceTags[0]] = instanceType
		}
	}
}

// GetInstanceType get all service type
func (s *InstancePool) GetInstanceType(tags string) []reflect.Type {
	if tags != "" {
		var types []reflect.Type
		if instance, ok := s.poolTags[tags]; ok {
			types = append(types, instance)
			return types
		}
		return types
	}
	return s.instanceTypeList
}

// CheckInstanceNameIfExist check contain name if exist
func (s *InstancePool) CheckInstanceNameIfExist(instanceName reflect.Type) bool {
	_, ok := s.poolMap[instanceName]
	return ok
}

// InstanceDISelfCheck  self check di request registered func exist or not
func (s *InstancePool) InstanceDISelfCheck(app Application) {
	for controllerName, pool := range s.poolMap {
		app.Logger().Infof("booting self DI checking instance %v ", controllerName)

		DiSelfCheck(controllerName, pool, app)
	}
	return

}

// GetInstance get service with di
func (s *InstancePool) GetInstance(instanceType reflect.Type, tctx Context, app Application, c *gin.Context) (interface{}, map[reflect.Type]interface{}) {
	pool, ok := s.poolMap[instanceType]
	if !ok {
		panic("unknown service name")
	}
	service := pool.Get()
	sharedInstance := DiAllFields(service, tctx, app, c)
	return service, sharedInstance
}

// InstanceMapping instance mapping
func (s *InstancePool) InstanceMapping() {}

// Release release service
func (s *InstancePool) Release(instance interface{}) {
	t := reflect.TypeOf(instance)
	syncpool, ok := s.poolMap[t]
	if !ok {
		return
	}
	DiFree(instance)
	syncpool.Put(instance)
}
