package service

import (
	"fmt"
	"metro_api/domain/model"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
)

var _ ResourceService = new(resourceServiceImpl)

func init() {
	trinitygo.RegisterInstance(resourceServiceImpl{}, "ResourceService")
}

// ResourceService  service interface
type ResourceService interface {
	GetResourceByID(id int64) (*model.Resource, error)
	GetResourceList(query string) (interface{}, error)
	CreateResource(*model.Resource) (*model.Resource, error)
	UpdateResourceByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteResourceByID(id int64, dVersion string) error
	MultiDeleteResourceByID([]modelutil.DeleteParam) error
}

type resourceServiceImpl struct {
	ResourceRepo repository.ResourceRepository `autowired:"true"  resource:"ResourceRepository"`
	Tctx         application.Context           `autowired:"true"`
}

func (s *resourceServiceImpl) GetResourceByID(id int64) (*model.Resource, error) {
	return s.ResourceRepo.GetResourceByID(id)
}
func (s *resourceServiceImpl) GetResourceList(query string) (interface{}, error) {
	res, isPaginationOff, err := s.ResourceRepo.GetResourceList(query)
	if err != nil {
		return nil, err
	}
	if isPaginationOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.ResourceRepo.GetResourceCount(query)
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

func (s *resourceServiceImpl) CreateResource(newResource *model.Resource) (*model.Resource, error) {
	return s.ResourceRepo.CreateResource(newResource)
}

func (s *resourceServiceImpl) UpdateResourceByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.ResourceRepo.UpdateResourceByID(id, dVersion, change)
}

func (s *resourceServiceImpl) DeleteResourceByID(id int64, dVersion string) error {
	return s.ResourceRepo.DeleteResourceByID(id, dVersion)
}

func (s *resourceServiceImpl) MultiDeleteResourceByID(deleteParam []modelutil.DeleteParam) error {
	for _, v := range deleteParam {
		if err := s.DeleteResourceByID(v.Key, v.DVersion); err != nil {
			return fmt.Errorf("Resource id %v deleted failed , %v", v.Key, err)
		}
	}
	return nil
}
