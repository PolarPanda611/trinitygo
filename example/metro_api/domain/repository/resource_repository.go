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
	_               ResourceRepository     = new(resourceRepositoryImpl)
	_resourceConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
		FilterBackend: []func(db *gorm.DB) *gorm.DB{},
		PageSize:      20,
		FilterList:    []string{"deleted_time__isnull", "code__ilike", "name__ilike", "description__ilike"},
		OrderByList:   []string{},
		SearchByList:  []string{},
		PreloadList: map[string]func(db *gorm.DB) *gorm.DB{
			"CreateUser": nil,
			"UpdateUser": nil,
			"DeleteUser": nil,
		},
		FilterCustomizeFunc: map[string]interface{}{},
	}
)

func init() {
	trinitygo.RegisterInstance(func() interface{} {
		repo := new(resourceRepositoryImpl)
		repo.queryHandler = queryutil.New(_resourceConfig)
		return repo
	}, "ResourceRepository")
}

// ResourceRepository resource repository
type ResourceRepository interface {
	GetResourceByID(id int64) (*model.Resource, error)
	GetResourceList(query string) ([]model.Resource, bool, error)
	CreateResource(*model.Resource) (*model.Resource, error)
	UpdateResourceByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteResourceByID(id int64, dVersion string) error
	GetResourceCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error)
}

type resourceRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *resourceRepositoryImpl) GetResourceByID(id int64) (*model.Resource, error) {
	var resource model.Resource
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&resource).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *resourceRepositoryImpl) GetResourceList(query string) ([]model.Resource, bool, error) {
	var resourceList []model.Resource
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&resourceList).Error; err != nil {
		return nil, false, err
	}
	return resourceList, r.queryHandler.IsPaginationOff(), nil
}

func (r *resourceRepositoryImpl) GetResourceCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.Resource{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *resourceRepositoryImpl) CreateResource(newResource *model.Resource) (*model.Resource, error) {
	if err := r.Tctx.DB().Create(newResource).Error; err != nil {
		return nil, err
	}
	return newResource, nil
}

func (r *resourceRepositoryImpl) UpdateResourceByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Table(r.Tctx.DB().NewScope(&model.Resource{}).TableName()).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}
	return nil
}
func (r *resourceRepositoryImpl) DeleteResourceByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.Resource{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}
	return nil
}
