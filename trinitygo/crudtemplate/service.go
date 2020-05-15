package crudtemplate

func init() {
	_templates["/domain/service/%v_service.go"] = genSrv()
}

func genSrv() string {
	return `
package service

import (
	"metro_api/domain/model"
	"strconv"

	"metro_api/domain/repository"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
)

var _ {{.ModelName}}Service = new({{.ModelNamePrivate}}ServiceImpl)

func init() {
	trinitygo.RegisterInstance({{.ModelNamePrivate}}ServiceImpl{}, "{{.ModelName}}Service")
}

// {{.ModelName}}Service  service interface
type {{.ModelName}}Service interface {
	Get{{.ModelName}}ByID(id int64) (*model.{{.ModelName}}, error)
	Get{{.ModelName}}List(query string) (interface{}, error)
	Create{{.ModelName}}(*model.{{.ModelName}}) (*model.{{.ModelName}}, error)
	Update{{.ModelName}}ByID(id int64, dVersion string, change map[string]interface{}) error
	Delete{{.ModelName}}ByID(id int64, dVersion string) error
}

type {{.ModelNamePrivate}}ServiceImpl struct {
	{{.ModelName}}Repo repository.{{.ModelName}}Repository ` + "`" + `autowired:"true"  resource:"{{.ModelName}}Repository"` + "`" + `
	Tctx     application.Context       ` + "`" + `autowired:"true"` + "`" + `
}

func (s *{{.ModelNamePrivate}}ServiceImpl) Get{{.ModelName}}ByID(id int64) (*model.{{.ModelName}}, error) {
	return s.{{.ModelName}}Repo.Get{{.ModelName}}ByID(id)
}
func (s *{{.ModelNamePrivate}}ServiceImpl) Get{{.ModelName}}List(query string) (interface{}, error) {
	res, err := s.{{.ModelName}}Repo.Get{{.ModelName}}List(query)
	if err != nil {
		return nil, err
	}
	IsOff, _ := strconv.ParseBool(s.Tctx.GinCtx().Query("PaginationOff"))
	if IsOff {
		return res, nil
	}
	count, currentPage, totalPage,pageSize, err := s.{{.ModelName}}Repo.Get{{.ModelName}}Count(query)
	if err != nil {
		return nil, err
	}
	resWithPagination := map[string]interface{}{
		"data":         res,
		"current_page": currentPage,
		"total_count":  count,
		"total_page":   totalPage,
		"page_size: 	pageSize,
	}
	return resWithPagination, nil
}

func (s *{{.ModelNamePrivate}}ServiceImpl) Create{{.ModelName}}(new{{.ModelName}} *model.{{.ModelName}}) (*model.{{.ModelName}}, error) {
	return s.{{.ModelName}}Repo.Create{{.ModelName}}(new{{.ModelName}})
}

func (s *{{.ModelNamePrivate}}ServiceImpl) Update{{.ModelName}}ByID(id int64, dVersion string, change map[string]interface{}) error {
	return s.{{.ModelName}}Repo.Update{{.ModelName}}ByID(id, dVersion, change)
}

func (s *{{.ModelNamePrivate}}ServiceImpl) Delete{{.ModelName}}ByID(id int64, dVersion string) error {
	return s.{{.ModelName}}Repo.Delete{{.ModelName}}ByID(id, dVersion)
}

	
	
	`
}
