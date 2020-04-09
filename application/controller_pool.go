package application

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// ControllerPool service pool
type ControllerPool struct {
	poolMap       map[string]*sync.Pool
	controllerMap []string
}

// NewControllerPool new pool with init map
func NewControllerPool() *ControllerPool {
	result := new(ControllerPool)
	result.poolMap = make(map[string]*sync.Pool)
	return result

}

// NewController add new service
func (s *ControllerPool) NewController(controllerType string, controllerPool *sync.Pool) {
	s.poolMap[controllerType] = controllerPool
	s.controllerMap = append(s.controllerMap, controllerType)
}

// GetControllerMap get controller map
func (s *ControllerPool) GetControllerMap() []string {
	return s.controllerMap
}

// GetController from pool
func (s *ControllerPool) GetController(controllerName string, tctx Context, app Application, c *gin.Context) (interface{}, []interface{}) {
	pool, ok := s.poolMap[controllerName]
	if !ok {
		panic("unknown controller name")
	}
	controller := pool.Get()
	toFreeContainer := DiFields(controller, tctx, app, c)
	return controller, toFreeContainer
}

// Release release controller to pool
func (s *ControllerPool) Release(controllerName string, controller interface{}) {
	pool, ok := s.poolMap[controllerName]
	if !ok {
		panic("unknown controller name")
	}
	DiFree(controller)
	pool.Put(controller)

}
