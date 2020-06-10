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

var _ GroupController = new(groupControllerImpl)

func init() {
	trinitygo.RegisterController("/v1/groups", groupControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetGroupByID"),
		application.NewRequestMapping(httputil.GET, "", "GetGroupList"),
		application.NewRequestMapping(httputil.POST, "", "CreateGroup"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateGroupByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteGroupByID"),
		application.NewRequestMapping(httputil.DELETE, "", "MultiDeleteGroupByID"),
	)
}

// GroupController group controller
type GroupController interface {
	GetGroupByID()
	GetGroupList()
	CreateGroup()
	UpdateGroupByID()
	DeleteGroupByID()
	MultiDeleteGroupByID()
}

type groupControllerImpl struct {
	GroupSrv service.GroupService `autowired:"true" resource:"GroupService"`
	Tctx     application.Context  `autowired:"true" transaction:"true"`
}

// GetGroupByID Method
// @Summary Get Group By ID
// @Description function for GroupController  to get Group By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/groups/{id} [get]
func (c *groupControllerImpl) GetGroupByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.GroupSrv.GetGroupByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// GetGroupList Method
// @Summary Get Group list By filter
// @Description function for GroupController  to get Group list By filter
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/groups [get]
func (c *groupControllerImpl) GetGroupList() {
	res, err := c.GroupSrv.GetGroupList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// CreateGroup Method
// @Summary Create Group
// @Description function for GroupController  to create Group
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/groups [post]
func (c *groupControllerImpl) CreateGroup() {
	var newGroup model.Group
	if err := c.Tctx.GinCtx().BindJSON(&newGroup); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.GroupSrv.CreateGroup(&newGroup)
	c.Tctx.httpResponseCreated(res, err)
	return
}

// UpdateGroupByID Method
// @Summary Modify Group
// @Description function for GroupController  to Modify Group
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/groups/{id} [patch]
func (c *groupControllerImpl) UpdateGroupByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.GroupSrv.UpdateGroupByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// DeleteGroupByID Method
// @Summary Delete Group
// @Description function for GroupController  to delete Group
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/groups/{id} [delete]
func (c *groupControllerImpl) DeleteGroupByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.GroupSrv.DeleteGroupByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// MultiDeleteGroupByID Method
// @Summary MultiDelete Group
// @Description function for GroupController  to MultiDelete Group
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/groups [delete]
func (c *groupControllerImpl) MultiDeleteGroupByID() {
	var deleteBody []modelutil.DeleteParam
	if err := c.Tctx.GinCtx().BindJSON(&deleteBody); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	err := c.GroupSrv.MultiDeleteGroupByID(deleteBody)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
