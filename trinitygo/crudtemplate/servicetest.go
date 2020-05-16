package crudtemplate

func init() {
	_templates["/domain/service/%v_service_test.go"] = genSrvTest()
}

func genSrvTest() string {
	return `
package service

import (
	"errors"
	"fmt"
	"{{.ProjectName}}/domain/model"
	"{{.ProjectName}}/domain/repository"
	"testing"

	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
	"github.com/PolarPanda611/trinitygo/testutil"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

var _ repository.{{.ModelName}}Repository = new({{.ModelName}}RepositoryMock)

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

var (
	// Init Mock Repository
	{{.ModelName}}Repo = new({{.ModelName}}RepositoryMock)
	// Init Mock Service
	{{.ModelName}}Srv = new({{.ModelNamePrivate}}ServiceImpl)
)

func init() {
	// Bind Mock Repo to Mock Service
	{{.ModelName}}Srv.{{.ModelName}}Repo = {{.ModelName}}Repo
}

// Test{{.ModelName}}ServiceGet{{.ModelName}}ByID test func Get{{.ModelName}}ByID for {{.ModelName}}Service
func Test{{.ModelName}}ServiceGet{{.ModelName}}ByID(t *testing.T) {
	// case 1
	{{.ModelName}}Repo.On("Get{{.ModelName}}ByID", int64(1)).Once().Return(&model.{{.ModelName}}{Code: "123"}, nil)
	testutil.Play(t, {{.ModelName}}Srv, "Get{{.ModelName}}ByID", int64(1)).Match(&model.{{.ModelName}}{Code: "123"}, nil)

	// case 2
	{{.ModelName}}Repo.On("Get{{.ModelName}}ByID", int64(2)).Once().Return(nil, gorm.ErrRecordNotFound)
	testutil.Play(t, {{.ModelName}}Srv, "Get{{.ModelName}}ByID", int64(2)).Match(nil, gorm.ErrRecordNotFound)
}

func Test{{.ModelName}}ServiceGet{{.ModelName}}List(t *testing.T) {
	// case 1
	{{.ModelName}}Repo.On("Get{{.ModelName}}List", "t").Once().Return(nil, false, errors.New("Get{{.ModelName}}List err"))
	testutil.Play(t, {{.ModelName}}Srv, "Get{{.ModelName}}List", "t").Match(nil, errors.New("Get{{.ModelName}}List err"))

	// case 2
	test{{.ModelName}}List := []model.{{.ModelName}}{
		model.{{.ModelName}}{
			Code: "1",
		},
		model.{{.ModelName}}{
			Code: "2",
		},
	}
	{{.ModelName}}Repo.On("Get{{.ModelName}}List", "t").Once().Return(test{{.ModelName}}List, true, nil)
	testutil.Play(t, {{.ModelName}}Srv, "Get{{.ModelName}}List", "t").Match(test{{.ModelName}}List, nil)

	// case 3
	{{.ModelName}}Repo.On("Get{{.ModelName}}List", "t").Once().Return(test{{.ModelName}}List, false, nil)
	{{.ModelName}}Repo.On("Get{{.ModelName}}Count", "t").Once().Return(0, 0, 0, 0, errors.New("Get{{.ModelName}}Count err"))
	testutil.Play(t, {{.ModelName}}Srv, "Get{{.ModelName}}List", "t").Match(nil, errors.New("Get{{.ModelName}}Count err"))

	// // case 4
	{{.ModelName}}Repo.On("Get{{.ModelName}}List", "t").Once().Return(test{{.ModelName}}List, false, nil)
	{{.ModelName}}Repo.On("Get{{.ModelName}}Count", "t").Once().Return(20, 1, 2, 10, nil)
	res := map[string]interface{}{
		"data":         test{{.ModelName}}List,
		"current_page": 1,
		"total_count":  20,
		"total_page":   2,
		"page_size":    10,
	}
	{{.ModelName}}Repo.On("Get{{.ModelName}}List", "t").Once().Return(test{{.ModelName}}List, false, nil)
	testutil.Play(t, {{.ModelName}}Srv, "Get{{.ModelName}}List", "t").Match(res, nil)
}

func Test{{.ModelName}}ServiceCreate{{.ModelName}}(t *testing.T) {
	// case 1
	{{.ModelName}}Repo.On("Create{{.ModelName}}", &model.{{.ModelName}}{Code: "123"}).Once().Return(&model.{{.ModelName}}{Code: "123"}, nil)
	testutil.Play(t, {{.ModelName}}Srv, "Create{{.ModelName}}", &model.{{.ModelName}}{Code: "123"}).Match(&model.{{.ModelName}}{Code: "123"}, nil)

	// case 2
	{{.ModelName}}Repo.On("Create{{.ModelName}}", &model.{{.ModelName}}{Code: "1234"}).Once().Return(nil, errors.New("Duplicate"))
	testutil.Play(t, {{.ModelName}}Srv, "Create{{.ModelName}}", &model.{{.ModelName}}{Code: "1234"}).Match(nil, errors.New("Duplicate"))
}

func Test{{.ModelName}}ServiceUpdate{{.ModelName}}ByID(t *testing.T) {
	// case 1
	{{.ModelName}}Repo.On("Update{{.ModelName}}ByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Once().Return(nil)
	testutil.Play(t, {{.ModelName}}Srv, "Update{{.ModelName}}ByID", int64(1), "v1", map[string]interface{}{"name": "123"}).Match(nil)

	// case 2
	{{.ModelName}}Repo.On("Update{{.ModelName}}ByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Once().Return(errors.New("Duplicate"))
	testutil.Play(t, {{.ModelName}}Srv, "Update{{.ModelName}}ByID", int64(1), "v1", map[string]interface{}{"name": "1234"}).Match(errors.New("Duplicate"))
}

func Test{{.ModelName}}ServiceDelete{{.ModelName}}ByID(t *testing.T) {
	// case 1
	{{.ModelName}}Repo.On("Delete{{.ModelName}}ByID", int64(1), "v1").Once().Return(nil)
	testutil.Play(t, {{.ModelName}}Srv, "Delete{{.ModelName}}ByID", int64(1), "v1").Match(nil)

	// case 2
	{{.ModelName}}Repo.On("Delete{{.ModelName}}ByID", int64(1), "v1").Once().Return(errors.New("Duplicate"))
	testutil.Play(t, {{.ModelName}}Srv, "Delete{{.ModelName}}ByID", int64(1), "v1").Match(errors.New("Duplicate"))
}

func Test{{.ModelName}}ServiceMultiDelete{{.ModelName}}ByID(t *testing.T) {
	// case 1
	deleteBody1 := []modelutil.DeleteParam{
		modelutil.DeleteParam{
			Key:      1,
			DVersion: "v1",
		},
		modelutil.DeleteParam{
			Key:      2,
			DVersion: "v1",
		},
	}
	{{.ModelName}}Repo.On("Delete{{.ModelName}}ByID", int64(1), "v1").Once().Return(nil)
	{{.ModelName}}Repo.On("Delete{{.ModelName}}ByID", int64(2), "v1").Once().Return(nil)
	testutil.Play(t, {{.ModelName}}Srv, "MultiDelete{{.ModelName}}ByID", deleteBody1).Match(nil)

	// case 2
	{{.ModelName}}Repo.On("Delete{{.ModelName}}ByID", int64(1), "v1").Once().Return(nil)
	{{.ModelName}}Repo.On("Delete{{.ModelName}}ByID", int64(2), "v1").Once().Return(errors.New("failed"))
	testutil.Play(t, {{.ModelName}}Srv, "MultiDelete{{.ModelName}}ByID", deleteBody1).Match(fmt.Errorf("e{{.ModelName}} id %v deleted failed , %v", int64(2), errors.New("failed")))
}

`
}
