package service

import (
	"reflect"
	"sync"
	"trinitygo"
	"trinitygo/application"
	"trinitygo/example/server/domain/repository"
)

func init() {
	trinitygo.BindService(reflect.TypeOf(&UserServiceImpl{}), &sync.Pool{
		New: func() interface{} {
			service := new(UserServiceImpl)
			return service
		},
	})
}

// UserService user service
type UserService interface {
	GetUserNameByID(name string) string
}

// UserServiceImpl user service iompl
type UserServiceImpl struct {
	UserRepo repository.UserRepo
	TContext application.Context
}

// GetUserNameByID method
func (s *UserServiceImpl) GetUserNameByID(name string) string {
	// fmt.Println("service runtime : ", s.TContext.GetRuntime())
	return s.UserRepo.Print()
}
