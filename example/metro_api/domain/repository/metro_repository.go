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
	_            MetroRepository        = new(metroRepositoryImpl)
	_metroConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
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
		repo := new(metroRepositoryImpl)
		repo.queryHandler = queryutil.New(_metroConfig)
		return repo
	}, "MetroRepository")
}

// MetroRepository metro repository
type MetroRepository interface {
	GetMetroByID(id int64) (*model.Metro, error)
	GetMetroList(query string) ([]model.Metro, bool, error)
	CreateMetro(*model.Metro) (*model.Metro, error)
	UpdateMetroByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteMetroByID(id int64, dVersion string) error
	GetMetroCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error)
}

type metroRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *metroRepositoryImpl) GetMetroByID(id int64) (*model.Metro, error) {
	var metro model.Metro
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&metro).Error; err != nil {
		return nil, err
	}
	return &metro, nil
}

func (r *metroRepositoryImpl) GetMetroList(query string) ([]model.Metro, bool, error) {
	var metroList []model.Metro
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&metroList).Error; err != nil {
		return nil, false, err
	}
	return metroList, r.queryHandler.IsPaginationOff(), nil
}

func (r *metroRepositoryImpl) GetMetroCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.Metro{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *metroRepositoryImpl) CreateMetro(newMetro *model.Metro) (*model.Metro, error) {
	if err := r.Tctx.DB().Create(newMetro).Error; err != nil {
		return nil, err
	}
	return newMetro, nil

}

func (r *metroRepositoryImpl) UpdateMetroByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Table(r.Tctx.DB().NewScope(&model.Metro{}).TableName()).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}

	return nil
}
func (r *metroRepositoryImpl) DeleteMetroByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.Metro{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}

	return nil
}
