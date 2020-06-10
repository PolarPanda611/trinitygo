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

var _ LineController = new(lineControllerImpl)

func init() {
	trinitygo.RegisterController("/v1/lines", lineControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetLineByID"),
		application.NewRequestMapping(httputil.GET, "", "GetLineList"),
		application.NewRequestMapping(httputil.POST, "", "CreateLine"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateLineByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteLineByID"),
		application.NewRequestMapping(httputil.DELETE, "", "MultiDeleteLineByID"),
	)
}

// LineController line controller
type LineController interface {
	GetLineByID()
	GetLineList()
	CreateLine()
	UpdateLineByID()
	DeleteLineByID()
	MultiDeleteLineByID()
}

type lineControllerImpl struct {
	LineSrv service.LineService `autowired:"true" resource:"LineService"`
	Tctx    application.Context `autowired:"true" transaction:"true"`
}

// GetLineByID Method
// @Summary Get Line By ID
// @Description function for controller to get Line By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/lines/{id} [get]
func (c *lineControllerImpl) GetLineByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.LineSrv.GetLineByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// GetLineList Method
// @Summary Get Line list
// @Description function for controller to get Line list
// @accept  json
// @Produce json
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/lines [get]
func (c *lineControllerImpl) GetLineList() {
	res, err := c.LineSrv.GetLineList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// CreateLine Method
// @Summary create Line list
// @Description function for controller to create Line list
// @accept  json
// @Produce json
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/lines [post]
func (c *lineControllerImpl) CreateLine() {
	var newLine model.Line
	if err := c.Tctx.GinCtx().BindJSON(&newLine); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.LineSrv.CreateLine(&newLine)
	c.Tctx.httpResponseCreated(res, err)
	return
}

// UpdateLineByID Method
// @Summary update Line by id
// @Description function for controller to create Line list
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/lines/{id} [patch]
func (c *lineControllerImpl) UpdateLineByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.LineSrv.UpdateLineByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// DeleteLineByID Method
// @Summary delete Line by id
// @Description function for controller to delete Line list
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/lines/{id} [delete]
func (c *lineControllerImpl) DeleteLineByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.LineSrv.DeleteLineByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// MultiDeleteLineByID Method
// @Summary multi delete Line by id
// @Description function for controller to multi delete Line list
// @accept  json
// @Produce json
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/lines [delete]
func (c *lineControllerImpl) MultiDeleteLineByID() {
	var deleteBody []modelutil.DeleteParam
	if err := c.Tctx.GinCtx().BindJSON(&deleteBody); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	err := c.LineSrv.MultiDeleteLineByID(deleteBody)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
