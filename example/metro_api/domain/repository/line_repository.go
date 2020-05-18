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
	_           LineRepository         = new(lineRepositoryImpl)
	_lineConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
		FilterBackend: []func(db *gorm.DB) *gorm.DB{},
		PageSize:      20,
		FilterList:    []string{"deleted_time__isnull", "code__ilike", "name__ilike", "description__ilike"},
		OrderByList:   []string{},
		SearchByList:  []string{"code", "name"},
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
		repo := new(lineRepositoryImpl)
		repo.queryHandler = queryutil.New(_lineConfig)
		return repo
	}, "LineRepository")
}

// LineRepository line repository
type LineRepository interface {
	GetLineByID(id int64) (*model.Line, error)
	GetLineList(query string) ([]model.Line, bool, error)
	CreateLine(*model.Line) (*model.Line, error)
	UpdateLineByID(id int64, dVersion string, change map[string]interface{}) error
	GetLineCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error)
	DeleteLineByID(id int64, dVersion string) error
}

type lineRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *lineRepositoryImpl) GetLineByID(id int64) (*model.Line, error) {
	var line model.Line
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&line).Error; err != nil {
		return nil, err
	}
	return &line, nil
}

func (r *lineRepositoryImpl) GetLineList(query string) ([]model.Line, bool, error) {
	var lineList []model.Line
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&lineList).Error; err != nil {
		return nil, false, err
	}
	return lineList, r.queryHandler.IsPaginationOff(), nil
}

func (r *lineRepositoryImpl) GetLineCount(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.Line{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *lineRepositoryImpl) CreateLine(newLine *model.Line) (*model.Line, error) {
	if err := r.Tctx.DB().Create(newLine).Error; err != nil {
		return nil, err
	}
	return newLine, nil

}

func (r *lineRepositoryImpl) UpdateLineByID(id int64, dVersion string, change map[string]interface{}) error {
	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Table(r.Tctx.DB().NewScope(&model.Line{}).TableName()).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}
	return nil
}

func (r *lineRepositoryImpl) DeleteLineByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.Line{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}
	return nil
}
