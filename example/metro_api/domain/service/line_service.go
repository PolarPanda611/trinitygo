package service

import (
	"errors"
	"fmt"
	"metro_api/domain/model"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
	"github.com/PolarPanda611/trinitygo/util"
)

var _ LineService = new(lineServiceImpl)

func init() {
	trinitygo.RegisterInstance(lineServiceImpl{}, "LineService")
}

// LineService  service interface
type LineService interface {
	GetLineByID(id int64) (*model.Line, error)
	GetLineList(query string) (interface{}, error)
	CreateLine(*model.Line) (*model.Line, error)
	UpdateLineByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteLineByID(id int64, dVersion string) error
	MultiDeleteLineByID([]modelutil.DeleteParam) error
}

type lineServiceImpl struct {
	LineRepo repository.LineRepository `autowired:"true"  resource:"LineRepository"`
	Tctx     application.Context       `autowired:"true"`
}

func (s *lineServiceImpl) GetLineByID(id int64) (*model.Line, error) {
	return s.LineRepo.GetLineByID(id)
}
func (s *lineServiceImpl) GetLineList(query string) (interface{}, error) {
	res, isPaginationOff, err := s.LineRepo.GetLineList(query)
	if err != nil {
		return nil, err
	}
	if isPaginationOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.LineRepo.GetLineCount(query)
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

func (s *lineServiceImpl) CreateLine(newLine *model.Line) (*model.Line, error) {
	v := util.ValueValidationImpl{}
	v.Load(newLine.Code, newLine.Name)
	if v.IfHasNilValue() {
		return nil, errors.New("Code , Name can not be empty")
	}
	return s.LineRepo.CreateLine(newLine)
}

func (s *lineServiceImpl) UpdateLineByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.LineRepo.UpdateLineByID(id, dVersion, change)
}

func (s *lineServiceImpl) DeleteLineByID(id int64, dVersion string) error {
	return s.LineRepo.DeleteLineByID(id, dVersion)
}

func (s *lineServiceImpl) MultiDeleteLineByID(deleteParam []modelutil.DeleteParam) error {
	for _, v := range deleteParam {
		if err := s.DeleteLineByID(v.Key, v.DVersion); err != nil {
			return fmt.Errorf("Line id %v deleted failed , %v", v.Key, err)
		}
	}
	return nil
}
