package repository

import (
	"metro_api/domain/model"

	"github.com/stretchr/testify/mock"
)

var _ AssetRepository = new(AssetRepositoryMock)

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
