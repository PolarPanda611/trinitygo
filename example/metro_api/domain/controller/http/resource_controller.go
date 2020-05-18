package http

import (
	"metro_api/domain/model"
	"strconv"

	"metro_api/domain/service"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
	"github.com/PolarPanda611/trinitygo/crud/util"
	"github.com/PolarPanda611/trinitygo/httputil"
)

var _ ResourceController = new(resourceControllerImpl)

func init() {
	trinitygo.RegisterController("/v1/resources", resourceControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetResourceByID"),
		application.NewRequestMapping(httputil.GET, "", "GetResourceList"),
		application.NewRequestMapping(httputil.POST, "", "CreateResource"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateResourceByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteResourceByID"),
		application.NewRequestMapping(httputil.DELETE, "", "MultiDeleteResourceByID"),
	)
}

// ResourceController resource controller
type ResourceController interface {
	GetResourceByID()
	GetResourceList()
	CreateResource()
	UpdateResourceByID()
	DeleteResourceByID()
	MultiDeleteResourceByID()
}

type resourceControllerImpl struct {
	ResourceSrv service.ResourceService `autowired:"true" resource:"ResourceService"`
	Tctx        application.Context     `autowired:"true" transaction:"true"`
}

// GetResourceByID Method
// @Summary Get Resource By ID
// @Description function for ResourceController  to get Resource By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/resources/{id} [get]
func (c *resourceControllerImpl) GetResourceByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.ResourceSrv.GetResourceByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// GetResourceList Method
// @Summary Get Resource list By filter
// @Description function for ResourceController  to get Resource list By filter
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/resources [get]
func (c *resourceControllerImpl) GetResourceList() {
	res, err := c.ResourceSrv.GetResourceList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// CreateResource Method
// @Summary Create Resource
// @Description function for ResourceController  to create Resource
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/resources [post]
func (c *resourceControllerImpl) CreateResource() {
	var newResource model.Resource
	if err := c.Tctx.GinCtx().BindJSON(&newResource); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.ResourceSrv.CreateResource(&newResource)
	c.Tctx.HTTPResponseCreated(res, err)
	return
}

// UpdateResourceByID Method
// @Summary Modify Resource
// @Description function for ResourceController  to Modify Resource
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/resources/{id} [patch]
func (c *resourceControllerImpl) UpdateResourceByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.ResourceSrv.UpdateResourceByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// DeleteResourceByID Method
// @Summary Delete Resource
// @Description function for ResourceController  to delete Resource
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/resources/{id} [delete]
func (c *resourceControllerImpl) DeleteResourceByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.ResourceSrv.DeleteResourceByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// MultiDeleteResourceByID Method
// @Summary MultiDelete Resource
// @Description function for ResourceController  to MultiDelete Resource
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/resources [delete]
func (c *resourceControllerImpl) MultiDeleteResourceByID() {
	var deleteBody []modelutil.DeleteParam
	if err := c.Tctx.GinCtx().BindJSON(&deleteBody); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	err := c.ResourceSrv.MultiDeleteResourceByID(deleteBody)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
