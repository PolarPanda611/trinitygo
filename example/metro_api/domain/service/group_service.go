package service

import (
	"fmt"
	"metro_api/domain/model"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
)

var _ GroupService = new(groupServiceImpl)

func init() {
	trinitygo.RegisterInstance(groupServiceImpl{}, "GroupService")
}

// GroupService  service interface
type GroupService interface {
	GetGroupByID(id int64) (*model.Group, error)
	GetGroupList(query string) (interface{}, error)
	CreateGroup(*model.Group) (*model.Group, error)
	UpdateGroupByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteGroupByID(id int64, dVersion string) error
	MultiDeleteGroupByID([]modelutil.DeleteParam) error
}

type groupServiceImpl struct {
	GroupRepo repository.GroupRepository `autowired:"true"  resource:"GroupRepository"`
	Tctx      application.Context        `autowired:"true"`
}

func (s *groupServiceImpl) GetGroupByID(id int64) (*model.Group, error) {
	return s.GroupRepo.GetGroupByID(id)
}
func (s *groupServiceImpl) GetGroupList(query string) (interface{}, error) {
	res, isPaginationOff, err := s.GroupRepo.GetGroupList(query)
	if err != nil {
		return nil, err
	}
	if isPaginationOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.GroupRepo.GetGroupCount(query)
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

func (s *groupServiceImpl) CreateGroup(newGroup *model.Group) (*model.Group, error) {
	return s.GroupRepo.CreateGroup(newGroup)
}

func (s *groupServiceImpl) UpdateGroupByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.GroupRepo.UpdateGroupByID(id, dVersion, change)
}

func (s *groupServiceImpl) DeleteGroupByID(id int64, dVersion string) error {
	return s.GroupRepo.DeleteGroupByID(id, dVersion)
}

func (s *groupServiceImpl) MultiDeleteGroupByID(deleteParam []modelutil.DeleteParam) error {
	for _, v := range deleteParam {
		if err := s.DeleteGroupByID(v.Key, v.DVersion); err != nil {
			return fmt.Errorf("Group id %v deleted failed , %v", v.Key, err)
		}
	}
	return nil
}
