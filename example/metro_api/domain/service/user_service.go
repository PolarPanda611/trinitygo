package service

import (
	"fmt"
	"metro_api/domain/model"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
)

var _ UserService = new(userServiceImpl)

func init() {
	trinitygo.RegisterInstance(userServiceImpl{}, "UserService")
}

// UserService  service interface
type UserService interface {
	GetUserByUserName(userName string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	GetUserList(query string) (interface{}, error)
	CreateUser(*model.User) (*model.User, error)
	UpdateUserByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteUserByID(id int64, dVersion string) error
	MultiDeleteUserByID([]modelutil.DeleteParam) error
}

type userServiceImpl struct {
	UserRepo repository.UserRepository `autowired:"true"  resource:"UserRepository"`
	Tctx     application.Context       `autowired:"true"`
}

func (s *userServiceImpl) GetUserByUserName(userName string) (*model.User, error) {
	return s.UserRepo.GetUserByUserName(userName)
}

func (s *userServiceImpl) GetUserByID(id int64) (*model.User, error) {
	return s.UserRepo.GetUserByID(id)
}
func (s *userServiceImpl) GetUserList(query string) (interface{}, error) {
	res, isPaginationOff, err := s.UserRepo.GetUserList(query)
	if err != nil {
		return nil, err
	}
	if isPaginationOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.UserRepo.GetUserCount(query)
	if err != nil {
		return nil, err
	}
	resWithPagination := map[string]interface{}{
		"data":       res,
		"current":    currentPage,
		"total":      count,
		"pageSize":   pageSize,
		"total_page": totalPage,
		"success":    true,
	}
	return resWithPagination, nil
}

func (s *userServiceImpl) CreateUser(newUser *model.User) (*model.User, error) {
	return s.UserRepo.CreateUser(newUser)
}

func (s *userServiceImpl) UpdateUserByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.UserRepo.UpdateUserByID(id, dVersion, change)
}

func (s *userServiceImpl) DeleteUserByID(id int64, dVersion string) error {
	return s.UserRepo.DeleteUserByID(id, dVersion)
}

func (s *userServiceImpl) MultiDeleteUserByID(deleteParam []modelutil.DeleteParam) error {
	for _, v := range deleteParam {
		if err := s.DeleteUserByID(v.Key, v.DVersion); err != nil {
			return fmt.Errorf("User id %v deleted failed , %v", v.Key, err)
		}
	}
	return nil
}
