package crudtemplate

func init() {
	_templates["/domain/repository/%v_repository.go"] = genRepo()
}

func genRepo() string {
	return `
package repository

import (
	"errors"
	"math"
	"{{.ProjectName}}/domain/model"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/queryutil"
	"github.com/jinzhu/gorm"
)

var (
	_           {{.ModelName}}Repository         = new({{.ModelNamePrivate}}RepositoryImpl)
	_{{.ModelNamePrivate}}Config *queryutil.QueryConfig = &queryutil.QueryConfig{
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
		repo := new({{.ModelNamePrivate}}RepositoryImpl)
		repo.queryHandler = queryutil.New(_{{.ModelNamePrivate}}Config)
		return repo
	}, "{{.ModelName}}Repository")
}

// {{.ModelName}}Repository {{.ModelNamePrivate}} repository
type {{.ModelName}}Repository interface {
	Get{{.ModelName}}ByID(id int64) (*model.{{.ModelName}}, error)
	Get{{.ModelName}}List(query string) ([]model.{{.ModelName}}, error)
	Create{{.ModelName}}(*model.{{.ModelName}}) (*model.{{.ModelName}}, error)
	Update{{.ModelName}}ByID(id int64, dVersion string, change map[string]interface{}) error
	Delete{{.ModelName}}ByID(id int64, dVersion string) error
	Get{{.ModelName}}Count(query string) (count uint, currentPage int, totalPage int,pageSize int,  err error)
}

type {{.ModelNamePrivate}}RepositoryImpl struct {
	Tctx         application.Context ` + "`" + `autowired:"true" ` + "`" + `
	queryHandler queryutil.QueryHandler
}

func (r *{{.ModelNamePrivate}}RepositoryImpl) Get{{.ModelName}}ByID(id int64) (*model.{{.ModelName}}, error) {
	var {{.ModelNamePrivate}} model.{{.ModelName}}
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ?", id).First(&{{.ModelNamePrivate}}).Error; err != nil {
		return nil, err
	}
	return &{{.ModelNamePrivate}}, nil
}

func (r *{{.ModelNamePrivate}}RepositoryImpl) Get{{.ModelName}}List(query string) ([]model.{{.ModelName}}, error) {
	var {{.ModelNamePrivate}}List []model.{{.ModelName}}
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Find(&{{.ModelNamePrivate}}List).Error; err != nil {
		return nil, err
	}
	return {{.ModelNamePrivate}}List, nil
}

func (r *{{.ModelNamePrivate}}RepositoryImpl) Get{{.ModelName}}Count(query string) (count uint, currentPage int, totalPage int,pageSize int,  err error) {
	if err := r.Tctx.DB().Scopes(
		r.queryHandler.HandleWithPagination(query)...,
	).Model(&model.{{.ModelName}}{}).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	return count, r.queryHandler.PageNum(), int(math.Ceil(float64(count) / float64(r.queryHandler.PageSize()))),r.queryHandler.PageSize(), nil
}

func (r *{{.ModelNamePrivate}}RepositoryImpl) Create{{.ModelName}}(new{{.ModelName}} *model.{{.ModelName}}) (*model.{{.ModelName}}, error) {
	if err := r.Tctx.DB().Create(new{{.ModelName}}).Error; err != nil {
		return nil, err
	}
	return new{{.ModelName}}, nil
}

func (r *{{.ModelNamePrivate}}RepositoryImpl) Update{{.ModelName}}ByID(id int64, dVersion string, change map[string]interface{}) error {

	updateQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Model(&model.{{.ModelName}}{}).Update(change)
	if err := updateQuery.Error; err != nil {
		return err
	}
	if updateQuery.RowsAffected != 1 {
		return errors.New("update affected zero lines , please refresh the data")
	}
	return nil
}
func (r *{{.ModelNamePrivate}}RepositoryImpl) Delete{{.ModelName}}ByID(id int64, dVersion string) error {
	deleteQuery := r.Tctx.DB().Scopes(
		r.queryHandler.HandleDBBackend()...,
	).Where("id = ? ", id).Where("d_version = ? ", dVersion).Delete(&model.{{.ModelName}}{})
	if err := deleteQuery.Error; err != nil {
		return err
	}
	if deleteQuery.RowsAffected != 1 {
		return errors.New("delete affected zero lines , please refresh the data")
	}
	return nil
}

	
	
`
}
