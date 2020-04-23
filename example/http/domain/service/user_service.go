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

var _ UserService = new(UserServiceImpl)

func init() {
	trinitygo.BindContainer(reflect.TypeOf(&UserServiceImpl{}), &sync.Pool{
		New: func() interface{} {
			service := new(UserServiceImpl)
			return service
		},
	},
		"UserService",
	)
}

//UserService user service
type UserService interface {
	GetUserByID(id int) (*object.User, error)
	GetUserList(query string) ([]object.User, error)
}

type UserServiceImpl struct {
	UserRepo repository.UserRepository `autowired:"true"`
	Tctx     application.Context       `autowired:"true"`
}

func (s *UserServiceImpl) GetUserByID(id int) (*object.User, error) {
	fmt.Println("service run ")
	return s.UserRepo.GetUserByID(id)
}

func (s *UserServiceImpl) GetUserList(query string) ([]object.User, error) {
	fmt.Println("service run ")
	return s.UserRepo.GetUserList(query)
}
