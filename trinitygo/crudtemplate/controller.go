package crudtemplate

func init() {
	_templates["/domain/controller/http/%v_controller.go"] = genController()
}

func genController() string {
	return `
package http

import (
	"metro_api/domain/model"
	"strconv"

	"metro_api/domain/service"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/crud/util"
	"github.com/PolarPanda611/trinitygo/httputil"
)

var _ {{.ModelName}}Controller = new({{.ModelNamePrivate}}ControllerImpl)

func init() {
	trinitygo.RegisterController("/{{.ModelNamePrivate}}s", {{.ModelNamePrivate}}ControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "Get{{.ModelName}}ByID"),
		application.NewRequestMapping(httputil.GET, "", "Get{{.ModelName}}List"),
		application.NewRequestMapping(httputil.POST, "", "Create{{.ModelName}}"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "Update{{.ModelName}}ByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "Delete{{.ModelName}}ByID"),
	)
}

// {{.ModelName}}Controller {{.ModelNamePrivate}} controller
type {{.ModelName}}Controller interface {
	Get{{.ModelName}}ByID()
	Get{{.ModelName}}List()
	Create{{.ModelName}}()
	Update{{.ModelName}}ByID()
	Delete{{.ModelName}}ByID()
}

type {{.ModelNamePrivate}}ControllerImpl struct {
	{{.ModelName}}Srv service.{{.ModelName}}Service ` + "`" + `autowired:"true" resource:"{{.ModelName}}Service"` + "`" + `
	Tctx    application.Context ` + "`" + `autowired:"true" transaction:"true"` + "`" + `
}

func (c *{{.ModelNamePrivate}}ControllerImpl) Get{{.ModelName}}ByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.{{.ModelName}}Srv.Get{{.ModelName}}ByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

func (c *{{.ModelNamePrivate}}ControllerImpl) Get{{.ModelName}}List() {
	res, err := c.{{.ModelName}}Srv.Get{{.ModelName}}List(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

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

	
	`
}
