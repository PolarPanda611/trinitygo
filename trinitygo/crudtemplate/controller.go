package crudtemplate

func init() {
	_templates["/domain/controller/http/%v_controller.go"] = genController()
}

func genController() string {
	return `
package http

import (
	"{{.ProjectName}}/domain/model"
	"{{.ProjectName}}/domain/service"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
	"github.com/PolarPanda611/trinitygo/httputil"
)

var _ {{.ModelName}}Controller = new({{.ModelNamePrivate}}ControllerImpl)

func init() {
	trinitygo.RegisterController("/{{.ModelNameToUnderscore}}s", {{.ModelNamePrivate}}ControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "Get{{.ModelName}}ByID"),
		application.NewRequestMapping(httputil.GET, "", "Get{{.ModelName}}List"),
		application.NewRequestMapping(httputil.POST, "", "Create{{.ModelName}}"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "Update{{.ModelName}}ByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "Delete{{.ModelName}}ByID"),
		application.NewRequestMapping(httputil.DELETE, "", "MultiDelete{{.ModelName}}ByID"),
	)
}

// {{.ModelName}}Controller {{.ModelNamePrivate}} controller
type {{.ModelName}}Controller interface {
	Get{{.ModelName}}ByID(args struct {
		ID int64 ` + "`" + `path_param:"id"` + "`" + `
	}) (*model.{{.ModelName}}, error)
	Get{{.ModelName}}List(args struct {
		Query string ` + "`" + `query_param:""` + "`" + `
	}) (interface{}, error)
	Create{{.ModelName}}(args struct {
		{{.ModelName}} model.{{.ModelName}} ` + "`" + `body_param:""` + "`" + `
	}) (*model.{{.ModelName}}, error)
	Update{{.ModelName}}ByID(args struct {
		ID       int64                  ` + "`" + `path_param:"id"` + "`" + `
		Change   map[string]interface{} ` + "`" + `body_param:""` + "`" + `
		DVersion string                 ` + "`" + `body_param:"d_version"` + "`" + `
	}) error 
	Delete{{.ModelName}}ByID(args struct {
		ID       int64  ` + "`" + `path_param:"id"` + "`" + `
		DVersion string ` + "`" + `body_param:"d_version"` + "`" + `
	}) error
	MultiDelete{{.ModelName}}ByID(args struct {
		DeleteParamList []modelutil.DeleteParam ` + "`" + `body_param:""` + "`" + `
	}) error 
}

type {{.ModelNamePrivate}}ControllerImpl struct {
	{{.ModelName}}Srv service.{{.ModelName}}Service ` + "`" + `autowired:"true" resource:"{{.ModelName}}Service"` + "`" + `
	Tctx    application.Context ` + "`" + `autowired:"true" transaction:"true"` + "`" + `
}

// Get{{.ModelName}}ByID Method
// @Summary Get {{.ModelName}} By ID
// @Description function for {{.ModelName}}Controller  to get {{.ModelName}} By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s/{id} [get]
func (c *{{.ModelNamePrivate}}ControllerImpl) Get{{.ModelName}}ByID(args struct {
	ID int64 ` + "`" + `path_param:"id"` + "`" + `
}) (*model.{{.ModelName}}, error) {
	return c.{{.ModelName}}Srv.Get{{.ModelName}}ByID(args.ID)
}

// Get{{.ModelName}}List Method
// @Summary Get {{.ModelName}} list By filter
// @Description function for {{.ModelName}}Controller  to get {{.ModelName}} list By filter
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s [get]
func (c *{{.ModelNamePrivate}}ControllerImpl) Get{{.ModelName}}List(args struct {
	Query string ` + "`" + `query_param:""` + "`" + `
}) (interface{}, error) {
	return c.{{.ModelName}}Srv.Get{{.ModelName}}List(args.Query)
}

// Create{{.ModelName}} Method
// @Summary Create {{.ModelName}} 
// @Description function for {{.ModelName}}Controller  to create {{.ModelName}} 
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s [post]
func (c *{{.ModelNamePrivate}}ControllerImpl) Create{{.ModelName}}(args struct {
	{{.ModelName}} model.{{.ModelName}} ` + "`" + `body_param:""` + "`" + `
})  (*model.{{.ModelName}}, error) {
	c.Tctx.HTTPStatus(201)
	return c.{{.ModelName}}Srv.Create{{.ModelName}}(&args.{{.ModelName}})
}

// Update{{.ModelName}}ByID Method
// @Summary Modify {{.ModelName}} 
// @Description function for {{.ModelName}}Controller  to Modify {{.ModelName}} 
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s/{id} [patch]
func (c *{{.ModelNamePrivate}}ControllerImpl) Update{{.ModelName}}ByID(args struct {
	ID       int64                  ` + "`" + `path_param:"id"` + "`" + `
	Change   map[string]interface{} ` + "`" + `body_param:""` + "`" + `
	DVersion string                 ` + "`" + `body_param:"d_version"` + "`" + `
}) error {
	return c.{{.ModelName}}Srv.Update{{.ModelName}}ByID(args.ID, args.DVersion, args.Change)
}

// Delete{{.ModelName}}ByID Method
// @Summary Delete {{.ModelName}} 
// @Description function for {{.ModelName}}Controller  to delete {{.ModelName}} 
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 204 {string} json "{"Status":204,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s/{id} [delete]
func (c *{{.ModelNamePrivate}}ControllerImpl) Delete{{.ModelName}}ByID(args struct {
	ID       int64  ` + "`" + `path_param:"id"` + "`" + `
	DVersion string ` + "`" + `body_param:"d_version"` + "`" + `
}) error {
	c.Tctx.HTTPStatus(204)
	return c.{{.ModelName}}Srv.Delete{{.ModelName}}ByID(args.ID, args.DVersion)
}

// MultiDelete{{.ModelName}}ByID Method
// @Summary MultiDelete {{.ModelName}} 
// @Description function for {{.ModelName}}Controller  to MultiDelete {{.ModelName}} 
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 204 {string} json "{"Status":204,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s [delete]
func (c *{{.ModelNamePrivate}}ControllerImpl) MultiDelete{{.ModelName}}ByID(args struct {
	DeleteParamList []modelutil.DeleteParam ` + "`" + `body_param:""` + "`" + `
}) error {
	c.Tctx.HTTPStatus(204)
	return c.{{.ModelName}}Srv.MultiDelete{{.ModelName}}ByID(args.DeleteParamList)
}
	
	`
}
