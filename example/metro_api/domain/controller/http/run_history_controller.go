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

var _ RunHistoryController = new(runHistoryControllerImpl)

func init() {
	trinitygo.RegisterController("/v1/run_histories", runHistoryControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetRunHistoryByID"),
		application.NewRequestMapping(httputil.GET, "", "GetRunHistoryList"),
		application.NewRequestMapping(httputil.POST, "", "CreateRunHistory"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateRunHistoryByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteRunHistoryByID"),
	)
}

// RunHistoryController runHistory controller
type RunHistoryController interface {
	GetRunHistoryByID()
	GetRunHistoryList()
	CreateRunHistory()
	UpdateRunHistoryByID()
	DeleteRunHistoryByID()
}

type runHistoryControllerImpl struct {
	RunHistorySrv service.RunHistoryService `autowired:"true" resource:"RunHistoryService"`
	Tctx          application.Context       `autowired:"true" transaction:"true"`
}

// GetRunHistoryByID Method
// @Summary Get RunHistory By ID
// @Description function for RunHistoryController  to get RunHistory By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/run_histories/{id} [get]
func (c *runHistoryControllerImpl) GetRunHistoryByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.RunHistorySrv.GetRunHistoryByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// GetRunHistoryList Method
// @Summary Get RunHistory list By filter
// @Description function for RunHistoryController  to get RunHistory list By filter
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/run_histories [get]
func (c *runHistoryControllerImpl) GetRunHistoryList() {
	res, err := c.RunHistorySrv.GetRunHistoryList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// CreateRunHistory Method
// @Summary Create RunHistory
// @Description function for RunHistoryController  to create RunHistory
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/run_histories [post]
func (c *runHistoryControllerImpl) CreateRunHistory() {
	var newRunHistory model.RunHistory
	if err := c.Tctx.GinCtx().BindJSON(&newRunHistory); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.RunHistorySrv.CreateRunHistory(&newRunHistory)
	c.Tctx.httpResponseCreated(res, err)
	return
}

// UpdateRunHistoryByID Method
// @Summary Modify RunHistory
// @Description function for RunHistoryController  to Modify RunHistory
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/run_histories/{id} [patch]
func (c *runHistoryControllerImpl) UpdateRunHistoryByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.RunHistorySrv.UpdateRunHistoryByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// DeleteRunHistoryByID Method
// @Summary Delete RunHistory
// @Description function for RunHistoryController  to delete RunHistory
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/run_histories/{id} [delete]
func (c *runHistoryControllerImpl) DeleteRunHistoryByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.RunHistorySrv.DeleteRunHistoryByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
