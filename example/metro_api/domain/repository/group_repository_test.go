package repository

import (
	"metro_api/domain/model"

	"github.com/stretchr/testify/mock"
)

var _ GroupRepository = new(GroupRepositoryMock)

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
