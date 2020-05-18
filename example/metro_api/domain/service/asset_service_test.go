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

var _ repository.AssetRepository = new(AssetRepositoryMock)

type AssetRepositoryMock struct {
	mock.Mock
}

func (r *AssetRepositoryMock) GetAssetByID(id int64) (*model.Asset, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *AssetRepositoryMock) GetAssetList(query string) ([]model.Asset, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Asset), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *AssetRepositoryMock) CreateAsset(asset *model.Asset) (*model.Asset, error) {
	args := r.Called(asset)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Asset), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *AssetRepositoryMock) UpdateAssetByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *AssetRepositoryMock) DeleteAssetByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *AssetRepositoryMock) GetAssetCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

var (
	// Init Mock Repository
	AssetRepo = new(AssetRepositoryMock)
	// Init Mock Service
	AssetSrv = new(assetServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	AssetSrv.AssetRepo = AssetRepo
}

// TestAssetServiceGetAssetByID test func GetAssetByID for AssetService
func TestAssetServiceGetAssetByID(t *testing.T) {
	// case 1
	AssetRepo.On("GetAssetByID", int64(1)).Once().Return(&model.Asset{Code: "123"}, nil)
	testutil.Play(t, AssetSrv, "GetAssetByID", int64(1)).Match(&model.Asset{Code: "123"}, nil)

	// case 2
	AssetRepo.On("GetAssetByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, AssetSrv, "GetAssetByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func TestAssetServiceGetAssetList(t *testing.T) {
	// case 1
	AssetRepo.On("GetAssetList", "t").Once().Return(nil, false, errors.New("GetAssetList err"))
	testutil.Play(t, AssetSrv, "GetAssetList", "t").Match(nil, errors.New("GetAssetList err"))

	// case 2
	testAssetList := []model.Asset{
		model.Asset{
			Code: "1",
		},
		model.Asset{
			Code: "2",
		},
	}
	AssetRepo.On("GetAssetList", "t").Once().Return(testAssetList, true, nil)
	testutil.Play(t, AssetSrv, "GetAssetList", "t").Match(testAssetList, nil)

	// case 3
	AssetRepo.On("GetAssetList", "t").Once().Return(testAssetList, false, nil)
	AssetRepo.On("GetAssetCount", "t").Once().Return(0, 0, 0, 0, errors.New("GetAssetCount err"))
	testutil.Play(t, AssetSrv, "GetAssetList", "t").Match(nil, errors.New("GetAssetCount err"))

	// // case 4
	AssetRepo.On("GetAssetList", "t").Once().Return(testAssetList, false, nil)
	AssetRepo.On("GetAssetCount", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":       testAssetList,
		"current":    1,
		"total":      20,
		"pageSize":   10,
		"total_page": 2,
		"success":    true,
	}
	AssetRepo.On("GetAssetList", "t").Once().Return(testAssetList, false, nil)
	testutil.Play(t, AssetSrv, "GetAssetList", "t").Match(res, nil)
}

func TestAssetServiceCreateAsset(t *testing.T) {
	// case 1
	AssetRepo.On("CreateAsset", &model.Asset{Code: "123"}).Once().Return(&model.Asset{Code: "123"}, nil)
	testutil.Play(t, AssetSrv, "CreateAsset", &model.Asset{Code: "123"}).Match(&model.Asset{Code: "123"}, nil)

	// case 2
	AssetRepo.On("CreateAsset", &model.Asset{Code: "1234"}).Once().Return(nil, errors.New("Duplicate"))
	testutil.Play(t, AssetSrv, "CreateAsset", &model.Asset{Code: "1234"}).Match(nil, errors.New("Duplicate"))
}

func TestAssetServiceUpdateAssetByID(t *testing.T) {
	// case 1
	AssetRepo.On("UpdateAssetByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, AssetSrv, "UpdateAssetByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	AssetRepo.On("UpdateAssetByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, AssetSrv, "UpdateAssetByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func TestAssetServiceDeleteAssetByID(t *testing.T) {
	// case 1
	AssetRepo.On("DeleteAssetByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, AssetSrv, "DeleteAssetByID", int64(1), "v1").Match(nil)

	// case 2
	AssetRepo.On("DeleteAssetByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, AssetSrv, "DeleteAssetByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func TestAssetServiceMultiDeleteAssetByID(t *testing.T) {
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
	AssetRepo.On("DeleteAssetByID", int64(1), "v1").Once().Return(nil)
	AssetRepo.On("DeleteAssetByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, AssetSrv, "MultiDeleteAssetByID", deleteBody1).Match(nil)

	// case 2
	AssetRepo.On("DeleteAssetByID", int64(1), "v1").Once().Return(nil)
	AssetRepo.On("DeleteAssetByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, AssetSrv, "MultiDeleteAssetByID", deleteBody1).Match(fmt.Errorf("Asset id %v deleted failed , %v", int64(2), errors.New("failed")))
}
