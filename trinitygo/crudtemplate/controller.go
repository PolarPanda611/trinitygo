package crudtemplate

func init() {
	_templates["/domain/controller/http/%v_controller.go"] = genController()
}

func genController() string {
	return `
package http

import (
	"{{.ProjectName}}/domain/model"
	"strconv"

	"{{.ProjectName}}/domain/service"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
	"github.com/PolarPanda611/trinitygo/crud/util"
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
	Get{{.ModelName}}ByID()
	Get{{.ModelName}}List()
	Create{{.ModelName}}()
	Update{{.ModelName}}ByID()
	Delete{{.ModelName}}ByID()
	MultiDelete{{.ModelName}}ByID()
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
func (c *{{.ModelNamePrivate}}ControllerImpl) Get{{.ModelName}}ByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.{{.ModelName}}Srv.Get{{.ModelName}}ByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
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
func (c *{{.ModelNamePrivate}}ControllerImpl) Get{{.ModelName}}List() {
	res, err := c.{{.ModelName}}Srv.Get{{.ModelName}}List(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
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
func (c *{{.ModelNamePrivate}}ControllerImpl) Create{{.ModelName}}() {
	var new{{.ModelName}} model.{{.ModelName}}
	if err := c.Tctx.GinCtx().BindJSON(&new{{.ModelName}}); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.{{.ModelName}}Srv.Create{{.ModelName}}(&new{{.ModelName}})
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// Update{{.ModelName}}ByID Method
// @Summary Modify {{.ModelName}} 
// @Description function for {{.ModelName}}Controller  to Modify {{.ModelName}} 
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s/{id} [patch]
func (c *{{.ModelNamePrivate}}ControllerImpl) Update{{.ModelName}}ByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.{{.ModelName}}Srv.Update{{.ModelName}}ByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// Delete{{.ModelName}}ByID Method
// @Summary Delete {{.ModelName}} 
// @Description function for {{.ModelName}}Controller  to delete {{.ModelName}} 
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s/{id} [delete]
func (c *{{.ModelNamePrivate}}ControllerImpl) Delete{{.ModelName}}ByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.{{.ModelName}}Srv.Delete{{.ModelName}}ByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// MultiDelete{{.ModelName}}ByID Method
// @Summary MultiDelete {{.ModelName}} 
// @Description function for {{.ModelName}}Controller  to MultiDelete {{.ModelName}} 
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /{{.ProjectName}}/{{.ModelNameToUnderscore}}s [delete]
func (c *{{.ModelNamePrivate}}ControllerImpl) MultiDelete{{.ModelName}}ByID() {
	var deleteBody []modelutil.DeleteParam
	if err := c.Tctx.GinCtx().BindJSON(&deleteBody); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	err := c.{{.ModelName}}Srv.MultiDelete{{.ModelName}}ByID(deleteBody)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
	
	`
}
