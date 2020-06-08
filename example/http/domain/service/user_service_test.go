package service

import (
	"errors"
	"fmt"
	"http/domain/model"
	"http/domain/repository"
	"testing"

	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
	"github.com/PolarPanda611/trinitygo/testutil"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

var _ repository.UserRepository = new(UserRepositoryMock)

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

var (
	// Init Mock Repository
	UserRepo = new(UserRepositoryMock)
	// Init Mock Service
	UserSrv = new(userServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	UserSrv.UserRepo = UserRepo
}

// TestUserServiceGetUserByID test func GetUserByID for UserService
func TestUserServiceGetUserByID(t *testing.T) {
	// case 1
	UserRepo.On("GetUserByID", int64(1)).Once().Return(&model.User{Code: "123"}, nil)
	testutil.Play(t, UserSrv, "GetUserByID", int64(1)).Match(&model.User{Code: "123"}, nil)

	// case 2
	UserRepo.On("GetUserByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, UserSrv, "GetUserByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func TestUserServiceGetUserList(t *testing.T) {
	// case 1
	UserRepo.On("GetUserList", "t").Once().Return(nil, false, errors.New("GetUserList err"))
	testutil.Play(t, UserSrv, "GetUserList", "t").Match(nil, errors.New("GetUserList err"))

	// case 2
	testUserList := []model.User{
		model.User{
			Code: "1",
		},
		model.User{
			Code: "2",
		},
	}
	UserRepo.On("GetUserList", "t").Once().Return(testUserList, true, nil)
	testutil.Play(t, UserSrv, "GetUserList", "t").Match(testUserList, nil)

	// case 3
	UserRepo.On("GetUserList", "t").Once().Return(testUserList, false, nil)
	UserRepo.On("GetUserCount", "t").Once().Return(0, 0, 0, 0, errors.New("GetUserCount err"))
	testutil.Play(t, UserSrv, "GetUserList", "t").Match(nil, errors.New("GetUserCount err"))

	// // case 4
	UserRepo.On("GetUserList", "t").Once().Return(testUserList, false, nil)
	UserRepo.On("GetUserCount", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":       testUserList,
		"current":    1,
		"total":      20,
		"pageSize":   10,
		"total_page": 2,
		"success":    true,
	}
	UserRepo.On("GetUserList", "t").Once().Return(testUserList, false, nil)
	testutil.Play(t, UserSrv, "GetUserList", "t").Match(res, nil)
}

func TestUserServiceCreateUser(t *testing.T) {
	// case 1
	UserRepo.On("CreateUser", &model.User{Code: "123"}).Once().Return(&model.User{Code: "123"}, nil)
	testutil.Play(t, UserSrv, "CreateUser", &model.User{Code: "123"}).Match(&model.User{Code: "123"}, nil)

	// case 2
	UserRepo.On("CreateUser", &model.User{Code: "1234"}).Once().Return(nil, errors.New("Duplicate"))
	testutil.Play(t, UserSrv, "CreateUser", &model.User{Code: "1234"}).Match(nil, errors.New("Duplicate"))
}

func TestUserServiceUpdateUserByID(t *testing.T) {
	// case 1
	UserRepo.On("UpdateUserByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, UserSrv, "UpdateUserByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	UserRepo.On("UpdateUserByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, UserSrv, "UpdateUserByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func TestUserServiceDeleteUserByID(t *testing.T) {
	// case 1
	UserRepo.On("DeleteUserByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, UserSrv, "DeleteUserByID", int64(1), "v1").Match(nil)

	// case 2
	UserRepo.On("DeleteUserByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, UserSrv, "DeleteUserByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func TestUserServiceMultiDeleteUserByID(t *testing.T) {
	// case 1
	deleteBody1 := []modelutil.DeleteParam{
		modelutil.DeleteParam{
			Key:      1,
			DVersion: "v1",
		},
		modelutil.DeleteParam{
			Key:      2,
			DVersion: "v1",
		},
	}
	UserRepo.On("DeleteUserByID", int64(1), "v1").Once().Return(nil)
	UserRepo.On("DeleteUserByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, UserSrv, "MultiDeleteUserByID", deleteBody1).Match(nil)

	// case 2
	UserRepo.On("DeleteUserByID", int64(1), "v1").Once().Return(nil)
	UserRepo.On("DeleteUserByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, UserSrv, "MultiDeleteUserByID", deleteBody1).Match(fmt.Errorf("User id %v deleted failed , %v", int64(2), errors.New("failed")))
}
