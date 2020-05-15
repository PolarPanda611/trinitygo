package service

import (
	"strconv"

	"github.com/PolarPanda611/trinitygo/example/http/domain/model"

	"github.com/PolarPanda611/trinitygo/example/http/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
)

var _ UserService = new(userServiceImpl)

func init() {
	trinitygo.RegisterInstance(userServiceImpl{}, "UserService")
}

// UserService  service interface
type UserService interface {
	//ServiceName
	//Method
	//Path
	GetUserByID(id int64) (*model.User, error)
	GetUserList(query string) (interface{}, error)
	CreateUser(*model.User) (*model.User, error)
	UpdateUserByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteUserByID(id int64, dVersion string) error
}

type userServiceImpl struct {
	UserRepo repository.UserRepository `autowired:"true"  resource:"UserRepository"`
	Tctx     application.Context       `autowired:"true"`
}

func (s *userServiceImpl) GetUserByID(id int64) (*model.User, error) {
	return s.UserRepo.GetUserByID(id)
}
func (s *userServiceImpl) GetUserList(query string) (interface{}, error) {
	res, err := s.UserRepo.GetUserList(query)
	if err != nil {
		return nil, err
	}
	IsOff, _ := strconv.ParseBool(s.Tctx.GinCtx().Query("PaginationOff"))
	if IsOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.UserRepo.GetUserCount(query)
	if err != nil {
		return nil, err
	}
	resWithPagination := map[string]interface{}{
		"data":         res,
		"current_page": currentPage,
		"total_count":  count,
		"total_page":   totalPage,
		"page_size":    pageSize,
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
