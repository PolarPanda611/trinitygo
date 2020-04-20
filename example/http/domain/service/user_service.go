package service

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/example/http/domain/object"
	"github.com/PolarPanda611/trinitygo/example/http/domain/repository"
)

var _ UserService = new(userServiceImpl)

func init() {
	trinitygo.BindContainer(reflect.TypeOf(&userServiceImpl{}), &sync.Pool{
		New: func() interface{} {
			service := new(userServiceImpl)
			return service
		},
	})
}

//UserService user service
type UserService interface {
	GetUserByID(id int) (*object.User, error)
	GetUserList(query string) ([]object.User, error)
}

type userServiceImpl struct {
	UserRepo repository.UserRepository
	TCtx     application.Context
}

func (s *userServiceImpl) GetUserByID(id int) (*object.User, error) {
	fmt.Println("service run ")
	return s.UserRepo.GetUserByID(id)
}

func (s *userServiceImpl) GetUserList(query string) ([]object.User, error) {
	fmt.Println("service run ")
	return s.UserRepo.GetUserList(query)
}
