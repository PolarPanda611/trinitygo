package crudtemplate

func init() {
	_templates["/domain/repository/%v_repository_test.go"] = genRepoTest()
}

func genRepoTest() string {
	return `
package repository

import (
	"{{.ProjectName}}/domain/model"

	"github.com/stretchr/testify/mock"
)

var _ {{.ModelName}}Repository = new({{.ModelName}}RepositoryMock)

type {{.ModelName}}RepositoryMock struct {
	mock.Mock
}

func (r *{{.ModelName}}RepositoryMock) Get{{.ModelName}}ByID(id int64) (*model.{{.ModelName}}, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.{{.ModelName}}), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *{{.ModelName}}RepositoryMock) Get{{.ModelName}}List(query string) ([]model.{{.ModelName}}, bool, error) {
	args := r.Called(query)
	if args.Get(0) != nil {
		return args.Get(0).([]model.{{.ModelName}}), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (r *{{.ModelName}}RepositoryMock) Create{{.ModelName}}({{.ModelNamePrivate}} *model.{{.ModelName}}) (*model.{{.ModelName}}, error) {
	args := r.Called({{.ModelNamePrivate}})
	if args.Get(0) != nil {
		return args.Get(0).(*model.{{.ModelName}}), args.Error(1)
	}
	return nil, args.Error(1)
}
func (r *{{.ModelName}}RepositoryMock) Update{{.ModelName}}ByID(id int64, dVersion string, change map[string]interface{}) error {
	args := r.Called(id, dVersion, change)
	return args.Error(0)
}
func (r *{{.ModelName}}RepositoryMock) Delete{{.ModelName}}ByID(id int64, dVersion string) error {
	args := r.Called(id, dVersion)
	return args.Error(0)
}
func (r *{{.ModelName}}RepositoryMock) Get{{.ModelName}}Count(query string) (count int, currentPage int, totalPage int, pageSize int, err error) {
	args := r.Called(query)
	return args.Get(0).(int), args.Get(1).(int), args.Get(2).(int), args.Get(3).(int), args.Error(4)
}

`
}
