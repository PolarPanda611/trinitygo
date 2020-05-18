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

var _ repository.ResourceRepository = new(ResourceRepositoryMock)

type ResourceRepositoryMock struct {
	mock.Mock
}

func (r *ResourceRepositoryMock) GetResourceByID(id int64) (*model.Resource, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Resource), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *ResourceRepositoryMock) GetResourceList(query string) ([]model.Resource, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Resource), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *ResourceRepositoryMock) CreateResource(resource *model.Resource) (*model.Resource, error) {
	args := r.Called(resource)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Resource), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *ResourceRepositoryMock) UpdateResourceByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *ResourceRepositoryMock) DeleteResourceByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *ResourceRepositoryMock) GetResourceCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

var (
	// Init Mock Repository
	ResourceRepo = new(ResourceRepositoryMock)
	// Init Mock Service
	ResourceSrv = new(resourceServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	ResourceSrv.ResourceRepo = ResourceRepo
}

// TestResourceServiceGetResourceByID test func GetResourceByID for ResourceService
func TestResourceServiceGetResourceByID(t *testing.T) {
	// case 1
	ResourceRepo.On("GetResourceByID", int64(1)).Once().Return(&model.Resource{Code: "123"}, nil)
	testutil.Play(t, ResourceSrv, "GetResourceByID", int64(1)).Match(&model.Resource{Code: "123"}, nil)

	// case 2
	ResourceRepo.On("GetResourceByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, ResourceSrv, "GetResourceByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func TestResourceServiceGetResourceList(t *testing.T) {
	// case 1
	ResourceRepo.On("GetResourceList", "t").Once().Return(nil, false, errors.New("GetResourceList err"))
	testutil.Play(t, ResourceSrv, "GetResourceList", "t").Match(nil, errors.New("GetResourceList err"))

	// case 2
	testResourceList := []model.Resource{
		model.Resource{
			Code: "1",
		},
		model.Resource{
			Code: "2",
		},
	}
	ResourceRepo.On("GetResourceList", "t").Once().Return(testResourceList, true, nil)
	testutil.Play(t, ResourceSrv, "GetResourceList", "t").Match(testResourceList, nil)

	// case 3
	ResourceRepo.On("GetResourceList", "t").Once().Return(testResourceList, false, nil)
	ResourceRepo.On("GetResourceCount", "t").Once().Return(0, 0, 0, 0, errors.New("GetResourceCount err"))
	testutil.Play(t, ResourceSrv, "GetResourceList", "t").Match(nil, errors.New("GetResourceCount err"))

	// // case 4
	ResourceRepo.On("GetResourceList", "t").Once().Return(testResourceList, false, nil)
	ResourceRepo.On("GetResourceCount", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":       testResourceList,
		"current":    1,
		"total":      20,
		"pageSize":   10,
		"total_page": 2,
		"success":    true,
	}
	ResourceRepo.On("GetResourceList", "t").Once().Return(testResourceList, false, nil)
	testutil.Play(t, ResourceSrv, "GetResourceList", "t").Match(res, nil)
}

func TestResourceServiceCreateResource(t *testing.T) {
	// case 1
	ResourceRepo.On("CreateResource", &model.Resource{Code: "123"}).Once().Return(&model.Resource{Code: "123"}, nil)
	testutil.Play(t, ResourceSrv, "CreateResource", &model.Resource{Code: "123"}).Match(&model.Resource{Code: "123"}, nil)

	// case 2
	ResourceRepo.On("CreateResource", &model.Resource{Code: "1234"}).Once().Return(nil, errors.New("Duplicate"))
	testutil.Play(t, ResourceSrv, "CreateResource", &model.Resource{Code: "1234"}).Match(nil, errors.New("Duplicate"))
}

func TestResourceServiceUpdateResourceByID(t *testing.T) {
	// case 1
	ResourceRepo.On("UpdateResourceByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, ResourceSrv, "UpdateResourceByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	ResourceRepo.On("UpdateResourceByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, ResourceSrv, "UpdateResourceByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func TestResourceServiceDeleteResourceByID(t *testing.T) {
	// case 1
	ResourceRepo.On("DeleteResourceByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, ResourceSrv, "DeleteResourceByID", int64(1), "v1").Match(nil)

	// case 2
	ResourceRepo.On("DeleteResourceByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, ResourceSrv, "DeleteResourceByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func TestResourceServiceMultiDeleteResourceByID(t *testing.T) {
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
	ResourceRepo.On("DeleteResourceByID", int64(1), "v1").Once().Return(nil)
	ResourceRepo.On("DeleteResourceByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, ResourceSrv, "MultiDeleteResourceByID", deleteBody1).Match(nil)

	// case 2
	ResourceRepo.On("DeleteResourceByID", int64(1), "v1").Once().Return(nil)
	ResourceRepo.On("DeleteResourceByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, ResourceSrv, "MultiDeleteResourceByID", deleteBody1).Match(fmt.Errorf("Resource id %v deleted failed , %v", int64(2), errors.New("failed")))
}
