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
	_                 RunHistoryRepository   = new(runHistoryRepositoryImpl)
	_runHistoryConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
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
		repo := new(runHistoryRepositoryImpl)
		repo.queryHandler = queryutil.New(_runHistoryConfig)
		return repo
	}, "RunHistoryRepository")
}

// RunHistoryRepository runHistory repository
type RunHistoryRepository interface {
	GetRunHistoryByID(id int64) (*model.RunHistory, error)
	GetRunHistoryList(query string) ([]model.RunHistory, bool, error)
	CreateRunHistory(*model.RunHistory) (*model.RunHistory, error)
	UpdateRunHistoryByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteRunHistoryByID(id int64, dVersion string) error
	GetRunHistoryCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error)
}

type runHistoryRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *runHistoryRepositoryImpl) GetRunHistoryByID(id int64) (*model.RunHistory, error) {
	var runHistory model.RunHistory
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&runHistory).Error; err != nil {
		return nil, err
	}
	return &runHistory, nil
}

func (r *runHistoryRepositoryImpl) GetRunHistoryList(query string) ([]model.RunHistory, bool, error) {
	var runHistoryList []model.RunHistory
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&runHistoryList).Error; err != nil {
		return nil, false, err
	}
	return runHistoryList, r.queryHandler.IsPaginationOff(), nil
}

func (r *runHistoryRepositoryImpl) GetRunHistoryCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.RunHistory{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *runHistoryRepositoryImpl) CreateRunHistory(newRunHistory *model.RunHistory) (*model.RunHistory, error) {
	if err := r.Tctx.DB().Create(newRunHistory).Error; err != nil {
		return nil, err
	}
	return newRunHistory, nil

}

func (r *runHistoryRepositoryImpl) UpdateRunHistoryByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Table(r.Tctx.DB().NewScope(&model.RunHistory{}).TableName()).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}

	return nil
}
func (r *runHistoryRepositoryImpl) DeleteRunHistoryByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.RunHistory{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}

	return nil
}
