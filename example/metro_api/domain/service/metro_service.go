package service

import (
	"fmt"
	"metro_api/domain/model"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
)

var _ MetroService = new(metroServiceImpl)

func init() {
	trinitygo.RegisterInstance(metroServiceImpl{}, "MetroService")
}

// MetroService  service interface
type MetroService interface {
	GetMetroByID(id int64) (*model.Metro, error)
	GetMetroList(query string) (interface{}, error)
	CreateMetro(*model.Metro) (*model.Metro, error)
	UpdateMetroByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteMetroByID(id int64, dVersion string) error
	MultiDeleteMetroByID([]modelutil.DeleteParam) error
}

type metroServiceImpl struct {
	MetroRepo repository.MetroRepository `autowired:"true"  resource:"MetroRepository"`
	Tctx      application.Context        `autowired:"true"`
}

func (s *metroServiceImpl) GetMetroByID(id int64) (*model.Metro, error) {
	return s.MetroRepo.GetMetroByID(id)
}
func (s *metroServiceImpl) GetMetroList(query string) (interface{}, error) {
	res, isPaginationOff, err := s.MetroRepo.GetMetroList(query)
	if err != nil {
		return nil, err
	}
	if isPaginationOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.MetroRepo.GetMetroCount(query)
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

func (s *metroServiceImpl) CreateMetro(newMetro *model.Metro) (*model.Metro, error) {
	return s.MetroRepo.CreateMetro(newMetro)
}

func (s *metroServiceImpl) UpdateMetroByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.MetroRepo.UpdateMetroByID(id, dVersion, change)
}

func (s *metroServiceImpl) DeleteMetroByID(id int64, dVersion string) error {
	return s.MetroRepo.DeleteMetroByID(id, dVersion)
}

func (s *metroServiceImpl) MultiDeleteMetroByID(deleteParam []modelutil.DeleteParam) error {
	for _, v := range deleteParam {
		if err := s.DeleteMetroByID(v.Key, v.DVersion); err != nil {
			return fmt.Errorf("Metro id %v deleted failed , %v", v.Key, err)
		}
	}
	return nil
}
