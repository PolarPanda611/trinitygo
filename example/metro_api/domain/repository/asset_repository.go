package repository

import (
	"errors"
	"math"
	"metro_api/domain/model"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/queryutil"
	"github.com/jinzhu/gorm"
)

var (
	_            AssetRepository        = new(assetRepositoryImpl)
	_assetConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
		FilterBackend:       []func(db *gorm.DB) *gorm.DB{},
		PageSize:            20,
		FilterList:          []string{},
		OrderByList:         []string{},
		SearchByList:        []string{},
		PreloadList:         map[string]func(db *gorm.DB) *gorm.DB{},
		FilterCustomizeFunc: map[string]interface{}{},
	}
)

func init() {
	trinitygo.RegisterInstance(func() interface{} {
		repo := new(assetRepositoryImpl)
		repo.queryHandler = queryutil.New(_assetConfig)
		return repo
	}, "AssetRepository")
}

// AssetRepository asset repository
type AssetRepository interface {
	GetAssetByID(id int64) (*model.Asset, error)
	GetAssetList(query string) ([]model.Asset, bool, error)
	CreateAsset(*model.Asset) (*model.Asset, error)
	UpdateAssetByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteAssetByID(id int64, dVersion string) error
	GetAssetCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error)
}

type assetRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *assetRepositoryImpl) GetAssetByID(id int64) (*model.Asset, error) {
	var asset model.Asset
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepositoryImpl) GetAssetList(query string) ([]model.Asset, bool, error) {
	var assetList []model.Asset
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&assetList).Error; err != nil {
		return nil, false, err
	}
	return assetList, r.queryHandler.IsPaginationOff(), nil
}

func (r *assetRepositoryImpl) GetAssetCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.Asset{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *assetRepositoryImpl) CreateAsset(newAsset *model.Asset) (*model.Asset, error) {
	if err := r.Tctx.DB().Create(newAsset).Error; err != nil {
		return nil, err
	}
	return newAsset, nil
}

func (r *assetRepositoryImpl) UpdateAssetByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Table(r.Tctx.DB().NewScope(&model.Asset{}).TableName()).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}
	return nil
}
func (r *assetRepositoryImpl) DeleteAssetByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.Asset{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}
	return nil
}
