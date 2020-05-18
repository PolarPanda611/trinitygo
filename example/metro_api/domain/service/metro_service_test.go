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

var _ repository.MetroRepository = new(MetroRepositoryMock)

type MetroRepositoryMock struct {
	mock.Mock
}

func (r *MetroRepositoryMock) GetMetroByID(id int64) (*model.Metro, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Metro), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *MetroRepositoryMock) GetMetroList(query string) ([]model.Metro, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Metro), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *MetroRepositoryMock) CreateMetro(metro *model.Metro) (*model.Metro, error) {
	args := r.Called(metro)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Metro), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *MetroRepositoryMock) UpdateMetroByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *MetroRepositoryMock) DeleteMetroByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *MetroRepositoryMock) GetMetroCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

var (
	// Init Mock Repository
	MetroRepo = new(MetroRepositoryMock)
	// Init Mock Service
	MetroSrv = new(metroServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	MetroSrv.MetroRepo = MetroRepo
}

// TestMetroServiceGetMetroByID test func GetMetroByID for MetroService
func TestMetroServiceGetMetroByID(t *testing.T) {
	// case 1
	MetroRepo.On("GetMetroByID", int64(1)).Once().Return(&model.Metro{Code: "123"}, nil)
	testutil.Play(t, MetroSrv, "GetMetroByID", int64(1)).Match(&model.Metro{Code: "123"}, nil)

	// case 2
	MetroRepo.On("GetMetroByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, MetroSrv, "GetMetroByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func TestMetroServiceGetMetroList(t *testing.T) {
	// case 1
	MetroRepo.On("GetMetroList", "t").Once().Return(nil, false, errors.New("GetMetroList err"))
	testutil.Play(t, MetroSrv, "GetMetroList", "t").Match(nil, errors.New("GetMetroList err"))

	// case 2
	testMetroList := []model.Metro{
		model.Metro{
			Code: "1",
		},
		model.Metro{
			Code: "2",
		},
	}
	MetroRepo.On("GetMetroList", "t").Once().Return(testMetroList, true, nil)
	testutil.Play(t, MetroSrv, "GetMetroList", "t").Match(testMetroList, nil)

	// case 3
	MetroRepo.On("GetMetroList", "t").Once().Return(testMetroList, false, nil)
	MetroRepo.On("GetMetroCount", "t").Once().Return(0, 0, 0, 0, errors.New("GetMetroCount err"))
	testutil.Play(t, MetroSrv, "GetMetroList", "t").Match(nil, errors.New("GetMetroCount err"))

	// // case 4
	MetroRepo.On("GetMetroList", "t").Once().Return(testMetroList, false, nil)
	MetroRepo.On("GetMetroCount", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":       testMetroList,
		"current":    1,
		"total":      20,
		"pageSize":   10,
		"total_page": 2,
		"success":    true,
	}
	MetroRepo.On("GetMetroList", "t").Once().Return(testMetroList, false, nil)
	testutil.Play(t, MetroSrv, "GetMetroList", "t").Match(res, nil)
}

func TestMetroServiceCreateMetro(t *testing.T) {
	// case 1
	MetroRepo.On("CreateMetro", &model.Metro{Code: "123"}).Once().Return(&model.Metro{Code: "123"}, nil)
	testutil.Play(t, MetroSrv, "CreateMetro", &model.Metro{Code: "123"}).Match(&model.Metro{Code: "123"}, nil)

	// case 2
	MetroRepo.On("CreateMetro", &model.Metro{Code: "1234"}).Once().Return(nil, errors.New("Duplicate"))
	testutil.Play(t, MetroSrv, "CreateMetro", &model.Metro{Code: "1234"}).Match(nil, errors.New("Duplicate"))
}

func TestMetroServiceUpdateMetroByID(t *testing.T) {
	// case 1
	MetroRepo.On("UpdateMetroByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, MetroSrv, "UpdateMetroByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	MetroRepo.On("UpdateMetroByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, MetroSrv, "UpdateMetroByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func TestMetroServiceDeleteMetroByID(t *testing.T) {
	// case 1
	MetroRepo.On("DeleteMetroByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, MetroSrv, "DeleteMetroByID", int64(1), "v1").Match(nil)

	// case 2
	MetroRepo.On("DeleteMetroByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, MetroSrv, "DeleteMetroByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func TestMetroServiceMultiDeleteMetroByID(t *testing.T) {
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
	MetroRepo.On("DeleteMetroByID", int64(1), "v1").Once().Return(nil)
	MetroRepo.On("DeleteMetroByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, MetroSrv, "MultiDeleteMetroByID", deleteBody1).Match(nil)

	// case 2
	MetroRepo.On("DeleteMetroByID", int64(1), "v1").Once().Return(nil)
	MetroRepo.On("DeleteMetroByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, MetroSrv, "MultiDeleteMetroByID", deleteBody1).Match(fmt.Errorf("Metro id %v deleted failed , %v", int64(2), errors.New("failed")))
}
