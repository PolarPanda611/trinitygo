package service

import (
	"errors"
	"fmt"
	"metro_api/domain/model"
	"metro_api/domain/repository"
	"testing"

	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
	"github.com/PolarPanda611/trinitygo/testutil"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

var _ repository.GroupRepository = new(GroupRepositoryMock)

type GroupRepositoryMock struct {
	mock.Mock
}

func (r *GroupRepositoryMock) GetGroupByID(id int64) (*model.Group, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Group), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *GroupRepositoryMock) GetGroupList(query string) ([]model.Group, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Group), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *GroupRepositoryMock) CreateGroup(group *model.Group) (*model.Group, error) {
	args := r.Called(group)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Group), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *GroupRepositoryMock) UpdateGroupByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *GroupRepositoryMock) DeleteGroupByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *GroupRepositoryMock) GetGroupCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

var (
	// Init Mock Repository
	GroupRepo = new(GroupRepositoryMock)
	// Init Mock Service
	GroupSrv = new(groupServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	GroupSrv.GroupRepo = GroupRepo
}

// TestGroupServiceGetGroupByID test func GetGroupByID for GroupService
func TestGroupServiceGetGroupByID(t *testing.T) {
	// case 1
	GroupRepo.On("GetGroupByID", int64(1)).Once().Return(&model.Group{Code: "123"}, nil)
	testutil.Play(t, GroupSrv, "GetGroupByID", int64(1)).Match(&model.Group{Code: "123"}, nil)

	// case 2
	GroupRepo.On("GetGroupByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, GroupSrv, "GetGroupByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func TestGroupServiceGetGroupList(t *testing.T) {
	// case 1
	GroupRepo.On("GetGroupList", "t").Once().Return(nil, false, errors.New("GetGroupList err"))
	testutil.Play(t, GroupSrv, "GetGroupList", "t").Match(nil, errors.New("GetGroupList err"))

	// case 2
	testGroupList := []model.Group{
		model.Group{
			Code: "1",
		},
		model.Group{
			Code: "2",
		},
	}
	GroupRepo.On("GetGroupList", "t").Once().Return(testGroupList, true, nil)
	testutil.Play(t, GroupSrv, "GetGroupList", "t").Match(testGroupList, nil)

	// case 3
	GroupRepo.On("GetGroupList", "t").Once().Return(testGroupList, false, nil)
	GroupRepo.On("GetGroupCount", "t").Once().Return(0, 0, 0, 0, errors.New("GetGroupCount err"))
	testutil.Play(t, GroupSrv, "GetGroupList", "t").Match(nil, errors.New("GetGroupCount err"))

	// // case 4
	GroupRepo.On("GetGroupList", "t").Once().Return(testGroupList, false, nil)
	GroupRepo.On("GetGroupCount", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":       testGroupList,
		"current":    1,
		"total":      20,
		"pageSize":   10,
		"total_page": 2,
		"success":    true,
	}
	GroupRepo.On("GetGroupList", "t").Once().Return(testGroupList, false, nil)
	testutil.Play(t, GroupSrv, "GetGroupList", "t").Match(res, nil)
}

func TestGroupServiceCreateGroup(t *testing.T) {
	// case 1
	GroupRepo.On("CreateGroup", &model.Group{Code: "123"}).Once().Return(&model.Group{Code: "123"}, nil)
	testutil.Play(t, GroupSrv, "CreateGroup", &model.Group{Code: "123"}).Match(&model.Group{Code: "123"}, nil)

	// case 2
	GroupRepo.On("CreateGroup", &model.Group{Code: "1234"}).Once().Return(nil, errors.New("Duplicate"))
	testutil.Play(t, GroupSrv, "CreateGroup", &model.Group{Code: "1234"}).Match(nil, errors.New("Duplicate"))
}

func TestGroupServiceUpdateGroupByID(t *testing.T) {
	// case 1
	GroupRepo.On("UpdateGroupByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, GroupSrv, "UpdateGroupByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	GroupRepo.On("UpdateGroupByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, GroupSrv, "UpdateGroupByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func TestGroupServiceDeleteGroupByID(t *testing.T) {
	// case 1
	GroupRepo.On("DeleteGroupByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, GroupSrv, "DeleteGroupByID", int64(1), "v1").Match(nil)

	// case 2
	GroupRepo.On("DeleteGroupByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, GroupSrv, "DeleteGroupByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func TestGroupServiceMultiDeleteGroupByID(t *testing.T) {
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
	GroupRepo.On("DeleteGroupByID", int64(1), "v1").Once().Return(nil)
	GroupRepo.On("DeleteGroupByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, GroupSrv, "MultiDeleteGroupByID", deleteBody1).Match(nil)

	// case 2
	GroupRepo.On("DeleteGroupByID", int64(1), "v1").Once().Return(nil)
	GroupRepo.On("DeleteGroupByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, GroupSrv, "MultiDeleteGroupByID", deleteBody1).Match(fmt.Errorf("Group id %v deleted failed , %v", int64(2), errors.New("failed")))
}
