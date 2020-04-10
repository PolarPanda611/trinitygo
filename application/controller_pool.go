package application

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// ControllerPool service pool
// if grpc string is the full method of method
// if http os the GET@/ping/:id
// need to filter controllerFuncMap to filter funcname
type ControllerPool struct {
	mu                sync.RWMutex
	poolMap           map[string]*sync.Pool
	controllerMap     []string
	controllerFuncMap map[string]string
}

// NewControllerPool new pool with init map
func NewControllerPool() *ControllerPool {
	result := new(ControllerPool)
	result.mu.Lock()
	defer result.mu.Unlock()
	result.poolMap = make(map[string]*sync.Pool)
	result.controllerFuncMap = make(map[string]string)
	return result

}

// NewController add new service
func (s *ControllerPool) NewController(controllerPool *sync.Pool, controllerType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.poolMap[controllerType] = controllerPool
	s.controllerMap = append(s.controllerMap, controllerType)
}

// NewControllerFunc register funcname for controllertype
func (s *ControllerPool) NewControllerFunc(controllerType string, funcName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.controllerFuncMap[controllerType] = funcName
}

// GetControllerMap get controller map
func (s *ControllerPool) GetControllerMap() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.controllerMap
}

// GetController from pool
func (s *ControllerPool) GetController(controllerName string, tctx Context, app Application, c *gin.Context) (interface{}, []interface{}) {
	s.mu.RLock()
	pool, ok := s.poolMap[controllerName]
	s.mu.RUnlock()
	if !ok {
		panic("unknown controller name")
	}
	controller := pool.Get()
	toFreeContainer := DiAllFields(controller, tctx, app, c)
	return controller, toFreeContainer
}

// GetControllerFuncName get controller func name
func (s *ControllerPool) GetControllerFuncName(controllerName string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.controllerFuncMap) == 0 {
		return "", false
	}
	funcName, ok := s.controllerFuncMap[controllerName]
	return funcName, ok

}

// Release release controller to pool
func (s *ControllerPool) Release(controllerName string, controller interface{}) {
	s.mu.RLock()
	pool, ok := s.poolMap[controllerName]
	s.mu.RUnlock()
	if !ok {
		panic("unknown controller name")
	}
	DiFree(controller)
	pool.Put(controller)

}
