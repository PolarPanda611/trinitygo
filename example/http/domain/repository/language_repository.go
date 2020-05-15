package repository

import (
	"errors"
	"math"

	"github.com/PolarPanda611/trinitygo/example/http/domain/model"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/queryutil"
	"github.com/jinzhu/gorm"
)

var (
	_               LanguageRepository     = new(languageRepositoryImpl)
	_languageConfig *queryutil.QueryConfig = &queryutil.QueryConfig{
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
		repo := new(languageRepositoryImpl)
		repo.queryHandler = queryutil.New(_languageConfig)
		return repo
	}, "LanguageRepository")
}

// LanguageRepository language repository
type LanguageRepository interface {
	GetLanguageByID(id int64) (*model.Language, error)
	GetLanguageList(query string) ([]model.Language, error)
	CreateLanguage(*model.Language) (*model.Language, error)
	UpdateLanguageByID(id int64, dVersion string, change map[string]interface{}) error
	DeleteLanguageByID(id int64, dVersion string) error
	GetLanguageCount(query string) (count uint, currentPage int, totalPage int, pageSize int, err error)
}

type languageRepositoryImpl struct {
	Tctx         application.Context `autowired:"true" `
	queryHandler queryutil.QueryHandler
}

func (r *languageRepositoryImpl) GetLanguageByID(id int64) (*model.Language, error) {
	var language model.Language
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&language).Error; err != nil {
		return nil, err
	}
	return &language, nil
}

func (r *languageRepositoryImpl) GetLanguageList(query string) ([]model.Language, error) {
	var languageList []model.Language
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&languageList).Error; err != nil {
		return nil, err
	}
	return languageList, nil
}

func (r *languageRepositoryImpl) GetLanguageCount(query string) (count uint, currentPage int, totalPage int, pageSize int, err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.Language{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))), r.queryHandler.PageSize(), nil
}

func (r *languageRepositoryImpl) CreateLanguage(newLanguage *model.Language) (*model.Language, error) {
	if err := r.Tctx.DB().Create(newLanguage).Error; err != nil {
		return nil, err
	}
	return newLanguage, nil

}

func (r *languageRepositoryImpl) UpdateLanguageByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Model(&model.Language{}).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}

	return nil
}
func (r *languageRepositoryImpl) DeleteLanguageByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.Language{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}

	return nil
}
