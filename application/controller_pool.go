package application

import (
	"fmt"
	"reflect"
	"sync"

	"log"

	"github.com/gin-gonic/gin"
)

// ControllerPool service pool
// if grpc string is the full method of method
// if http os the GET@/ping/:id
// need to filter controllerFuncMap to filter funcname
type ControllerPool struct {
	mu                   sync.RWMutex
	poolMap              map[string]*sync.Pool
	controllerMap        []string
	controllerFuncMap    map[string]string
	controllerValidators map[string][]Validator
}

// NewControllerPool new pool with init map
func NewControllerPool() *ControllerPool {
	result := new(ControllerPool)
	result.mu.Lock()
	defer result.mu.Unlock()
	result.poolMap = make(map[string]*sync.Pool)
	result.controllerFuncMap = make(map[string]string)
	result.controllerValidators = make(map[string][]Validator)
	return result

}

// NewController add new service
func (s *ControllerPool) NewController(controllerType string, controllerPool *sync.Pool) {
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

// NewControllerValidators register funcname for controllertype
func (s *ControllerPool) NewControllerValidators(controllerType string, validator ...Validator) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.controllerValidators[controllerType] = validator
}

// GetControllerMap get controller map
func (s *ControllerPool) GetControllerMap() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.controllerMap
}

// ControllSelfCheck self check http request registered func exist or not
func (s *ControllerPool) ControllSelfCheck(controllerName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pool, poolExist := s.poolMap[controllerName]
	if !poolExist {
		log.Fatalf("controller %v not registered , self check failed", controllerName)
	}
	funcName, funcExist := s.controllerFuncMap[controllerName]
	if funcName == "" || !funcExist {
		// func not exist
		return false
	}
	controller := pool.Get()
	defer pool.Put(controller)
	_, funcImpled := reflect.TypeOf(controller).MethodByName(funcName)
	if !funcImpled {
		log.Fatalf("func %v not registered on controller %v , self check failed", funcName, controllerName)
	}
	return true
}

// GetController from pool
func (s *ControllerPool) GetController(controllerName string, tctx Context, app Application, c *gin.Context) (interface{}, []interface{}) {
	s.mu.RLock()
	pool, ok := s.poolMap[controllerName]
	s.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("unknown controller name : %v", controllerName))
	}
	controller := pool.Get()
	toFreeContainer := DiAllFields(controller, tctx, app, c)
	return controller, toFreeContainer
}

// GetControllerFuncName get controller func name
func (s *ControllerPool) GetControllerFuncName(controllerName string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.controllerFuncMap) == 0 {
		return "", false
	}
	funcName, ok := s.controllerFuncMap[controllerName]
	return funcName, ok

}

// GetControllerValidators get controller func name
func (s *ControllerPool) GetControllerValidators(controllerName string) []Validator {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.controllerValidators) == 0 {
		return nil
	}
	validators, ok := s.controllerValidators[controllerName]
	if !ok {
		return nil
	}
	return validators

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
