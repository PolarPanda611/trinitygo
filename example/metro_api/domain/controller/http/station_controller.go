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

var _ StationController = new(stationControllerImpl)

func init() {
	trinitygo.RegisterController("/v1/stations", stationControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetStationByID"),
		application.NewRequestMapping(httputil.GET, "", "GetStationList"),
		application.NewRequestMapping(httputil.POST, "", "CreateStation"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateStationByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteStationByID"),
	)
}

// StationController station controller
type StationController interface {
	GetStationByID()
	GetStationList()
	CreateStation()
	UpdateStationByID()
	DeleteStationByID()
}

type stationControllerImpl struct {
	StationSrv service.StationService `autowired:"true" resource:"StationService"`
	Tctx       application.Context    `autowired:"true" transaction:"true"`
}

// GetStationByID Method
// @Summary Get Station By ID
// @Description function for StationController  to get Station By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/stations/{id} [get]
func (c *stationControllerImpl) GetStationByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.StationSrv.GetStationByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// GetStationList Method
// @Summary Get Station list By filter
// @Description function for StationController  to get Station list By filter
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/stations [get]
func (c *stationControllerImpl) GetStationList() {
	res, err := c.StationSrv.GetStationList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// CreateStation Method
// @Summary Create Station
// @Description function for StationController  to create Station
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/stations [post]
func (c *stationControllerImpl) CreateStation() {
	var newStation model.Station
	if err := c.Tctx.GinCtx().BindJSON(&newStation); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.StationSrv.CreateStation(&newStation)
	c.Tctx.httpResponseCreated(res, err)
	return
}

// UpdateStationByID Method
// @Summary Modify Station
// @Description function for StationController  to Modify Station
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/stations/{id} [patch]
func (c *stationControllerImpl) UpdateStationByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.StationSrv.UpdateStationByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// DeleteStationByID Method
// @Summary Delete Station
// @Description function for StationController  to delete Station
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/stations/{id} [delete]
func (c *stationControllerImpl) DeleteStationByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.StationSrv.DeleteStationByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
