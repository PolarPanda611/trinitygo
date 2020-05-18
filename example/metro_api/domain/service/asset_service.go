package service

import (
	"fmt"
	"metro_api/domain/model"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
)

var _ AssetService = new(assetServiceImpl)

func init() {
	trinitygo.RegisterInstance(assetServiceImpl{}, "AssetService")
}

// AssetService  service interface
type AssetService interface {
	GetAssetByID(id int64) (*model.Asset, error)
	GetAssetList(query string) (interface{}, error)
	CreateAsset(*model.Asset) (*model.Asset, error)
	UpdateAssetByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteAssetByID(id int64, dVersion string) error
	MultiDeleteAssetByID([]modelutil.DeleteParam) error
}

type assetServiceImpl struct {
	AssetRepo repository.AssetRepository `autowired:"true"  resource:"AssetRepository"`
	Tctx      application.Context        `autowired:"true"`
}

func (s *assetServiceImpl) GetAssetByID(id int64) (*model.Asset, error) {
	return s.AssetRepo.GetAssetByID(id)
}
func (s *assetServiceImpl) GetAssetList(query string) (interface{}, error) {
	res, isPaginationOff, err := s.AssetRepo.GetAssetList(query)
	if err != nil {
		return nil, err
	}
	if isPaginationOff {
		return res, nil
	}
	count, currentPage, totalPage, pageSize, err := s.AssetRepo.GetAssetCount(query)
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

func (s *assetServiceImpl) CreateAsset(newAsset *model.Asset) (*model.Asset, error) {
	return s.AssetRepo.CreateAsset(newAsset)
}

func (s *assetServiceImpl) UpdateAssetByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.AssetRepo.UpdateAssetByID(id, dVersion, change)
}

func (s *assetServiceImpl) DeleteAssetByID(id int64, dVersion string) error {
	return s.AssetRepo.DeleteAssetByID(id, dVersion)
}

func (s *assetServiceImpl) MultiDeleteAssetByID(deleteParam []modelutil.DeleteParam) error {
	for _, v := range deleteParam {
		if err := s.DeleteAssetByID(v.Key, v.DVersion); err != nil {
			return fmt.Errorf("Asset id %v deleted failed , %v", v.Key, err)
		}
	}
	return nil
}
