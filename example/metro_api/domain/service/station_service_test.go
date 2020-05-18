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

var _ repository.StationRepository = new(StationRepositoryMock)

type StationRepositoryMock struct {
	mock.Mock
}

func (r *StationRepositoryMock) GetNextSeq() (string, error) {
	args := r.Called()
	return args.String(0), args.Error(1)
}
func (r *StationRepositoryMock) GetStationByID(id int64) (*model.Station, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Station), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *StationRepositoryMock) GetStationList(query string) ([]model.Station, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Station), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *StationRepositoryMock) CreateStation(station *model.Station) (*model.Station, error) {
	args := r.Called(station)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Station), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *StationRepositoryMock) UpdateStationByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *StationRepositoryMock) DeleteStationByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *StationRepositoryMock) GetStationCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

var (
	// Init Mock Repository
	StationRepo = new(StationRepositoryMock)
	// Init Mock Service
	StationSrv = new(stationServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	StationSrv.StationRepo = StationRepo
}

// TestStationServiceGetNextStationCode  test func getNextStationCode for StationService
func TestStationServiceGetNextStationCode(t *testing.T) {
	// case 1
	StationRepo.On("GetNextSeq").Once().Return("1", nil)
	res1, res2 := StationSrv.getNextStationCode()
	testutil.PlayUnexported(t, res1, res2).Match("S00001", nil)

	// case 2
	StationRepo.On("GetNextSeq").Once().Return("", errors.New("fuck"))
	res1, res2 = StationSrv.getNextStationCode()
	testutil.PlayUnexported(t, res1, res2).Match("", errors.New("fuck"))
}

// TestStationServiceGetStationByID test func GetStationByID for StationService
func TestStationServiceGetStationByID(t *testing.T) {
	// case 1
	StationRepo.On("GetStationByID", int64(1)).Once().Return(&model.Station{Code: "123"}, nil)
	testutil.Play(t, StationSrv, "GetStationByID", int64(1)).Match(&model.Station{Code: "123"}, nil)

	// case 2
	StationRepo.On("GetStationByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, StationSrv, "GetStationByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func TestStationServiceGetStationList(t *testing.T) {
	// case 1
	StationRepo.On("GetStationList", "t").Once().Return(nil, false, errors.New("GetStationList err"))
	testutil.Play(t, StationSrv, "GetStationList", "t").Match(nil, errors.New("GetStationList err"))

	// case 2
	testStationList := []model.Station{
		model.Station{
			Code: "1",
		},
		model.Station{
			Code: "2",
		},
	}
	StationRepo.On("GetStationList", "t").Once().Return(testStationList, true, nil)
	testutil.Play(t, StationSrv, "GetStationList", "t").Match(testStationList, nil)

	// case 3
	StationRepo.On("GetStationList", "t").Once().Return(testStationList, false, nil)
	StationRepo.On("GetStationCount", "t").Once().Return(0, 0, 0, 0, errors.New("GetStationCount err"))
	testutil.Play(t, StationSrv, "GetStationList", "t").Match(nil, errors.New("GetStationCount err"))

	// // case 4
	StationRepo.On("GetStationList", "t").Once().Return(testStationList, false, nil)
	StationRepo.On("GetStationCount", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":       testStationList,
		"current":    1,
		"total":      20,
		"pageSize":   10,
		"total_page": 2,
		"success":    true,
	}
	StationRepo.On("GetStationList", "t").Once().Return(testStationList, false, nil)
	testutil.Play(t, StationSrv, "GetStationList", "t").Match(res, nil)
}

func TestStationServiceCreateStation(t *testing.T) {
	// case 1
	StationRepo.On("GetNextSeq").Once().Return("1", nil)
	StationRepo.On("CreateStation", &model.Station{Code: "S00001", Name: "123"}).Once().Return(&model.Station{Code: "S00001", Name: "123"}, nil)
	testutil.Play(t, StationSrv, "CreateStation", &model.Station{Name: "123"}).Match(&model.Station{Code: "S00001", Name: "123"}, nil)

	// case 2
	StationRepo.On("GetNextSeq").Once().Return("1", nil)
	StationRepo.On("CreateStation", &model.Station{Code: "S00001", Name: "1234"}).Once().Return(nil, errors.New("Duplicate"))
	testutil.Play(t, StationSrv, "CreateStation", &model.Station{Name: "1234"}).Match(nil, errors.New("Duplicate"))
}

func TestStationServiceUpdateStationByID(t *testing.T) {
	// case 1
	StationRepo.On("UpdateStationByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, StationSrv, "UpdateStationByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	StationRepo.On("UpdateStationByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, StationSrv, "UpdateStationByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func TestStationServiceDeleteStationByID(t *testing.T) {
	// case 1
	StationRepo.On("DeleteStationByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, StationSrv, "DeleteStationByID", int64(1), "v1").Match(nil)

	// case 2
	StationRepo.On("DeleteStationByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, StationSrv, "DeleteStationByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func TestStationServiceMultiDeleteStationByID(t *testing.T) {
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
	StationRepo.On("DeleteStationByID", int64(1), "v1").Once().Return(nil)
	StationRepo.On("DeleteStationByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, StationSrv, "MultiDeleteStationByID", deleteBody1).Match(nil)

	// case 2
	StationRepo.On("DeleteStationByID", int64(1), "v1").Once().Return(nil)
	StationRepo.On("DeleteStationByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, StationSrv, "MultiDeleteStationByID", deleteBody1).Match(fmt.Errorf("Station id %v deleted failed , %v", int64(2), errors.New("failed")))
}
