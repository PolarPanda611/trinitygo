package application

import (
	"fmt"
	"reflect"
	"sync"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/kataras/golog"
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

// ControllerFuncSelfCheck self check http request registered func exist or not
func (s *ControllerPool) ControllerFuncSelfCheck(isLog bool, logger *golog.Logger) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for controllerName, pool := range s.poolMap {
		funcName, funcExist := s.controllerFuncMap[controllerName]
		if funcName == "" || !funcExist {
			// func not exist
			logger.Fatalf("booting self func checking controller %v , no func registered , self check failed ...", controllerName)
		}
		controller := pool.Get()
		defer pool.Put(controller)
		_, funcImpled := reflect.TypeOf(controller).MethodByName(funcName)
		if !funcImpled {
			log.Fatalf("booting self func checking controller %v , func %v not registered , self check failed ...", controllerName, funcName)
		}
		if isLog {
			logger.Infof("booting self func checking controller %v , func %v checked ", controllerName, funcName)
		}
	}

	return
}

// ControllerDISelfCheck  self check di request registered func exist or not
func (s *ControllerPool) ControllerDISelfCheck(app Application) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for controllerName, pool := range s.poolMap {
		app.Logger().Infof("booting self DI checking controller %v ", controllerName)

		DiSelfCheck(controllerName, pool, app)
	}
	return

}

// GetController from pool
func (s *ControllerPool) GetController(controllerName string, tctx Context, app Application, c *gin.Context) (interface{}, map[reflect.Type]interface{}) {
	s.mu.RLock()
	pool, ok := s.poolMap[controllerName]
	s.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("unknown controller name : %v", controllerName))
	}
	controller := pool.Get()
	sharedInstance := DiAllFields(controller, tctx, app, c, true)
	return controller, sharedInstance
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
