package repository

import (
	"metro_api/domain/model"

	"github.com/stretchr/testify/mock"
)

var _ ResourceRepository = new(ResourceRepositoryMock)

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
