package service

import (
	"fmt"
	"metro_api/domain/model"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
)

var _ RunHistoryService = new(runHistoryServiceImpl)

func init() {
	trinitygo.RegisterInstance(runHistoryServiceImpl{}, "RunHistoryService")
}

// RunHistoryService  service interface
type RunHistoryService interface {
	GetRunHistoryByID(id int64) (*model.RunHistory, error)
	GetRunHistoryList(query string) (interface{}, error)
	CreateRunHistory(*model.RunHistory) (*model.RunHistory, error)
	UpdateRunHistoryByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteRunHistoryByID(id int64, dVersion string) error
	MultiDeleteRunHistoryByID([]modelutil.DeleteParam) error
}

type runHistoryServiceImpl struct {
	RunHistoryRepo repository.RunHistoryRepository `autowired:"true"  resource:"RunHistoryRepository"`
	Tctx           application.Context             `autowired:"true"`
}

func (s *runHistoryServiceImpl) GetRunHistoryByID(id int64) (*model.RunHistory, error) {
	return s.RunHistoryRepo.GetRunHistoryByID(id)
}
func (s *runHistoryServiceImpl) GetRunHistoryList(query string) (interface{}, error) {
	res, isPaginationOff, err := s.RunHistoryRepo.GetRunHistoryList(query)
	if err != nil {
		return nil, err
	}
	if isPaginationOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.RunHistoryRepo.GetRunHistoryCount(query)
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

func (s *runHistoryServiceImpl) CreateRunHistory(newRunHistory *model.RunHistory) (*model.RunHistory, error) {
	return s.RunHistoryRepo.CreateRunHistory(newRunHistory)
}

func (s *runHistoryServiceImpl) UpdateRunHistoryByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.RunHistoryRepo.UpdateRunHistoryByID(id, dVersion, change)
}

func (s *runHistoryServiceImpl) DeleteRunHistoryByID(id int64, dVersion string) error {
	return s.RunHistoryRepo.DeleteRunHistoryByID(id, dVersion)
}

func (s *runHistoryServiceImpl) MultiDeleteRunHistoryByID(deleteParam []modelutil.DeleteParam) error {
	for _, v := range deleteParam {
		if err := s.DeleteRunHistoryByID(v.Key, v.DVersion); err != nil {
			return fmt.Errorf("RunHistory id %v deleted failed , %v", v.Key, err)
		}
	}
	return nil
}
