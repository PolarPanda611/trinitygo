package service

import (
	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/example/server/domain/repository"
)

var _ UserService = new(UserServiceImpl)

func init() {
	trinitygo.BindContainer(UserServiceImpl{})
}

// UserService user service
type UserService interface {
	GetUserNameByID(name string) string
}

// UserServiceImpl user service iompl
type UserServiceImpl struct {
	UserRepo repository.UserRepo `autowired:"true"`
	TContext application.Context `autowired:"true"`
}

// GetUserNameByID method
func (s *UserServiceImpl) GetUserNameByID(name string) string {
	// fmt.Println("service runtime : ", s.TContext.Runtime())
	return s.UserRepo.Print()
}
