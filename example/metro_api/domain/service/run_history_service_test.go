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

var _ repository.RunHistoryRepository = new(RunHistoryRepositoryMock)

type RunHistoryRepositoryMock struct {
	mock.Mock
}

func (r *RunHistoryRepositoryMock) GetRunHistoryByID(id int64) (*model.RunHistory, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.RunHistory), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *RunHistoryRepositoryMock) GetRunHistoryList(query string) ([]model.RunHistory, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.RunHistory), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *RunHistoryRepositoryMock) CreateRunHistory(runHistory *model.RunHistory) (*model.RunHistory, error) {
	args := r.Called(runHistory)
	if args.Get(0) != nil {
		return args.Get(0).(*model.RunHistory), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *RunHistoryRepositoryMock) UpdateRunHistoryByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *RunHistoryRepositoryMock) DeleteRunHistoryByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *RunHistoryRepositoryMock) GetRunHistoryCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

var (
	// Init Mock Repository
	RunHistoryRepo = new(RunHistoryRepositoryMock)
	// Init Mock Service
	RunHistorySrv = new(runHistoryServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	RunHistorySrv.RunHistoryRepo = RunHistoryRepo
}

// TestRunHistoryServiceGetRunHistoryByID test func GetRunHistoryByID for RunHistoryService
func TestRunHistoryServiceGetRunHistoryByID(t *testing.T) {
	// case 1
	RunHistoryRepo.On("GetRunHistoryByID", int64(1)).Once().Return(&model.RunHistory{Code: "123"}, nil)
	testutil.Play(t, RunHistorySrv, "GetRunHistoryByID", int64(1)).Match(&model.RunHistory{Code: "123"}, nil)

	// case 2
	RunHistoryRepo.On("GetRunHistoryByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, RunHistorySrv, "GetRunHistoryByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func TestRunHistoryServiceGetRunHistoryList(t *testing.T) {
	// case 1
	RunHistoryRepo.On("GetRunHistoryList", "t").Once().Return(nil, false, errors.New("GetRunHistoryList err"))
	testutil.Play(t, RunHistorySrv, "GetRunHistoryList", "t").Match(nil, errors.New("GetRunHistoryList err"))

	// case 2
	testRunHistoryList := []model.RunHistory{
		model.RunHistory{
			Code: "1",
		},
		model.RunHistory{
			Code: "2",
		},
	}
	RunHistoryRepo.On("GetRunHistoryList", "t").Once().Return(testRunHistoryList, true, nil)
	testutil.Play(t, RunHistorySrv, "GetRunHistoryList", "t").Match(testRunHistoryList, nil)

	// case 3
	RunHistoryRepo.On("GetRunHistoryList", "t").Once().Return(testRunHistoryList, false, nil)
	RunHistoryRepo.On("GetRunHistoryCount", "t").Once().Return(0, 0, 0, 0, errors.New("GetRunHistoryCount err"))
	testutil.Play(t, RunHistorySrv, "GetRunHistoryList", "t").Match(nil, errors.New("GetRunHistoryCount err"))

	// // case 4
	RunHistoryRepo.On("GetRunHistoryList", "t").Once().Return(testRunHistoryList, false, nil)
	RunHistoryRepo.On("GetRunHistoryCount", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":       testRunHistoryList,
		"current":    1,
		"total":      20,
		"pageSize":   10,
		"total_page": 2,
		"success":    true,
	}
	RunHistoryRepo.On("GetRunHistoryList", "t").Once().Return(testRunHistoryList, false, nil)
	testutil.Play(t, RunHistorySrv, "GetRunHistoryList", "t").Match(res, nil)
}

func TestRunHistoryServiceCreateRunHistory(t *testing.T) {
	// case 1
	RunHistoryRepo.On("CreateRunHistory", &model.RunHistory{Code: "123"}).Once().Return(&model.RunHistory{Code: "123"}, nil)
	testutil.Play(t, RunHistorySrv, "CreateRunHistory", &model.RunHistory{Code: "123"}).Match(&model.RunHistory{Code: "123"}, nil)

	// case 2
	RunHistoryRepo.On("CreateRunHistory", &model.RunHistory{Code: "1234"}).Once().Return(nil, errors.New("Duplicate"))
	testutil.Play(t, RunHistorySrv, "CreateRunHistory", &model.RunHistory{Code: "1234"}).Match(nil, errors.New("Duplicate"))
}

func TestRunHistoryServiceUpdateRunHistoryByID(t *testing.T) {
	// case 1
	RunHistoryRepo.On("UpdateRunHistoryByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, RunHistorySrv, "UpdateRunHistoryByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	RunHistoryRepo.On("UpdateRunHistoryByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, RunHistorySrv, "UpdateRunHistoryByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func TestRunHistoryServiceDeleteRunHistoryByID(t *testing.T) {
	// case 1
	RunHistoryRepo.On("DeleteRunHistoryByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, RunHistorySrv, "DeleteRunHistoryByID", int64(1), "v1").Match(nil)

	// case 2
	RunHistoryRepo.On("DeleteRunHistoryByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, RunHistorySrv, "DeleteRunHistoryByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func TestRunHistoryServiceMultiDeleteRunHistoryByID(t *testing.T) {
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
	RunHistoryRepo.On("DeleteRunHistoryByID", int64(1), "v1").Once().Return(nil)
	RunHistoryRepo.On("DeleteRunHistoryByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, RunHistorySrv, "MultiDeleteRunHistoryByID", deleteBody1).Match(nil)

	// case 2
	RunHistoryRepo.On("DeleteRunHistoryByID", int64(1), "v1").Once().Return(nil)
	RunHistoryRepo.On("DeleteRunHistoryByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, RunHistorySrv, "MultiDeleteRunHistoryByID", deleteBody1).Match(fmt.Errorf("RunHistory id %v deleted failed , %v", int64(2), errors.New("failed")))
}
