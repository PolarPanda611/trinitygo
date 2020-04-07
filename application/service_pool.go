package application

import (
	"reflect"
	"sync"
)

// ServicePool service pool
type ServicePool struct {
	poolMap     map[reflect.Type]*sync.Pool
	serviceType []reflect.Type
	serviceMap  sync.Map
}

// NewServicePool new pool with init map
func NewServicePool() *ServicePool {
	result := new(ServicePool)
	result.poolMap = make(map[reflect.Type]*sync.Pool)
	return result

}

// NewService add new service
func (s *ServicePool) NewService(srvType reflect.Type, srvPool *sync.Pool) {
	s.poolMap[srvType] = srvPool
	s.serviceType = append(s.serviceType, srvType)
}

// GetServiceType get all service type
func (s *ServicePool) GetServiceType() []reflect.Type {
	return s.serviceType
}

// GetService get service with di
func (s *ServicePool) GetService(serviceType reflect.Type, tctx Context, app Application) (interface{}, []interface{}) {
	pool, ok := s.poolMap[serviceType]
	if !ok {
		panic("unknown service name")
	}
	service := pool.Get()
	toFreeRepository := DiServicePool(service, tctx, app)
	return service, toFreeRepository
}

// Release release service
func (s *ServicePool) Release(service interface{}) {
	t := reflect.TypeOf(service)
	syncpool, ok := s.poolMap[t]
	if !ok {
		return
	}
	DiFree(service)
	syncpool.Put(service)
}
