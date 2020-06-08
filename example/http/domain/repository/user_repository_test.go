package repository

import (
	"http/domain/model"

	"github.com/stretchr/testify/mock"
)

var _ UserRepository = new(UserRepositoryMock)

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) GetUserByID(id int64) (*model.User, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *UserRepositoryMock) GetUserList(query string) ([]model.User, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.User), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *UserRepositoryMock) CreateUser(user *model.User) (*model.User, error) {
	args := r.Called(user)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *UserRepositoryMock) UpdateUserByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *UserRepositoryMock) DeleteUserByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *UserRepositoryMock) GetUserCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}
