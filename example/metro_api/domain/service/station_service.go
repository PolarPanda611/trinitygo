package service

import (
	"fmt"
	"metro_api/domain/model"
	"metro_api/infra/util"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
)

var _ StationService = new(stationServiceImpl)

func init() {
	trinitygo.RegisterInstance(stationServiceImpl{}, "StationService")
}

// StationService  service interface
type StationService interface {
	getNextStationCode() (string, error)
	GetStationByID(id int64) (*model.Station, error)
	GetStationList(query string) (interface{}, error)
	CreateStation(*model.Station) (*model.Station, error)
	UpdateStationByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteStationByID(id int64, dVersion string) error
	MultiDeleteStationByID([]modelutil.DeleteParam) error
}

type stationServiceImpl struct {
	StationRepo repository.StationRepository `autowired:"true"  resource:"StationRepository"`
	Tctx        application.Context          `autowired:"true"`
}

func (s *stationServiceImpl) getNextStationCode() (string, error) {
	nextSeq, err := s.StationRepo.GetNextSeq()
	if err != nil {
		return "", err
	}
	return util.GetCodeUtil("S", "00000", nextSeq, 5), nil
}

func (s *stationServiceImpl) GetStationByID(id int64) (*model.Station, error) {
	return s.StationRepo.GetStationByID(id)
}
func (s *stationServiceImpl) GetStationList(query string) (interface{}, error) {
	res, isPaginationOff, err := s.StationRepo.GetStationList(query)
	if err != nil {
		return nil, err
	}
	if isPaginationOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.StationRepo.GetStationCount(query)
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

func (s *stationServiceImpl) CreateStation(newStation *model.Station) (*model.Station, error) {
	nextCode, err := s.getNextStationCode()
	if err != nil {
		return nil, err
	}
	newStation.Code = nextCode
	return s.StationRepo.CreateStation(newStation)
}

func (s *stationServiceImpl) UpdateStationByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.StationRepo.UpdateStationByID(id, dVersion, change)
}

func (s *stationServiceImpl) DeleteStationByID(id int64, dVersion string) error {
	return s.StationRepo.DeleteStationByID(id, dVersion)
}

func (s *stationServiceImpl) MultiDeleteStationByID(deleteParam []modelutil.DeleteParam) error {
	for _, v := range deleteParam {
		if err := s.DeleteStationByID(v.Key, v.DVersion); err != nil {
			return fmt.Errorf("Station id %v deleted failed , %v", v.Key, err)
		}
	}
	return nil
}
