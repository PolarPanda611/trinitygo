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

var _ repository.LineRepository = new(LineRepositoryMock)

type LineRepositoryMock struct {
	mock.Mock
}

func (r *LineRepositoryMock) GetLineByID(id int64) (*model.Line, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Line), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *LineRepositoryMock) GetLineList(query string) ([]model.Line, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Line), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *LineRepositoryMock) CreateLine(line *model.Line) (*model.Line, error) {
	args := r.Called(line)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Line), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *LineRepositoryMock) UpdateLineByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *LineRepositoryMock) DeleteLineByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *LineRepositoryMock) GetLineCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

var (
	// Init Mock Repository
	LineRepo = new(LineRepositoryMock)
	// Init Mock Service
	LineSrv = new(lineServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	LineSrv.LineRepo = LineRepo
}

// TestLineServiceGetLineByID test func GetLineByID for LineService
func TestLineServiceGetLineByID(t *testing.T) {
	// case 1
	LineRepo.On("GetLineByID", int64(1)).Once().Return(&model.Line{Code: "123"}, nil)
	testutil.Play(t, LineSrv, "GetLineByID", int64(1)).Match(&model.Line{Code: "123"}, nil)

	// case 2
	LineRepo.On("GetLineByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, LineSrv, "GetLineByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func TestLineServiceGetLineList(t *testing.T) {
	// case 1
	LineRepo.On("GetLineList", "t").Once().Return(nil, false, errors.New("GetLineList err"))
	testutil.Play(t, LineSrv, "GetLineList", "t").Match(nil, errors.New("GetLineList err"))

	// case 2
	testLineList := []model.Line{
		model.Line{
			Code: "1",
		},
		model.Line{
			Code: "2",
		},
	}
	LineRepo.On("GetLineList", "t").Once().Return(testLineList, true, nil)
	testutil.Play(t, LineSrv, "GetLineList", "t").Match(testLineList, nil)

	// case 3
	LineRepo.On("GetLineList", "t").Once().Return(testLineList, false, nil)
	LineRepo.On("GetLineCount", "t").Once().Return(0, 0, 0, 0, errors.New("GetLineCount err"))
	testutil.Play(t, LineSrv, "GetLineList", "t").Match(nil, errors.New("GetLineCount err"))

	// // case 4
	LineRepo.On("GetLineList", "t").Once().Return(testLineList, false, nil)
	LineRepo.On("GetLineCount", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":       testLineList,
		"current":    1,
		"total":      20,
		"pageSize":   10,
		"total_page": 2,
		"success":    true,
	}
	LineRepo.On("GetLineList", "t").Once().Return(testLineList, false, nil)
	testutil.Play(t, LineSrv, "GetLineList", "t").Match(res, nil)
}

func TestLineServiceCreateLine(t *testing.T) {
	// case 1
	LineRepo.On("CreateLine", &model.Line{Code: "123", Name: "234"}).Once().Return(&model.Line{Code: "123", Name: "234"}, nil)
	testutil.Play(t, LineSrv, "CreateLine", &model.Line{Code: "123", Name: "234"}).Match(&model.Line{Code: "123", Name: "234"}, nil)

	// case 2
	testutil.Play(t, LineSrv, "CreateLine", &model.Line{Code: "1234"}).Match(nil, errors.New("Code , Name can not be empty"))

	// case 3
	LineRepo.On("CreateLine", &model.Line{Code: "1234", Name: "234"}).Once().Return(nil, errors.New("duplicate"))
	testutil.Play(t, LineSrv, "CreateLine", &model.Line{Code: "1234", Name: "234"}).Match(nil, errors.New("duplicate"))
}

func TestLineServiceUpdateLineByID(t *testing.T) {
	// case 1
	LineRepo.On("UpdateLineByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, LineSrv, "UpdateLineByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	LineRepo.On("UpdateLineByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, LineSrv, "UpdateLineByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func TestLineServiceDeleteLineByID(t *testing.T) {
	// case 1
	LineRepo.On("DeleteLineByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, LineSrv, "DeleteLineByID", int64(1), "v1").Match(nil)

	// case 2
	LineRepo.On("DeleteLineByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, LineSrv, "DeleteLineByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func TestLineServiceMultiDeleteLineByID(t *testing.T) {
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
	LineRepo.On("DeleteLineByID", int64(1), "v1").Once().Return(nil)
	LineRepo.On("DeleteLineByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, LineSrv, "MultiDeleteLineByID", deleteBody1).Match(nil)

	// case 2
	LineRepo.On("DeleteLineByID", int64(1), "v1").Once().Return(nil)
	LineRepo.On("DeleteLineByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, LineSrv, "MultiDeleteLineByID", deleteBody1).Match(fmt.Errorf("Line id %v deleted failed , %v", int64(2), errors.New("failed")))
}
