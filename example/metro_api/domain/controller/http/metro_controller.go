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

var _ MetroController = new(metroControllerImpl)

func init() {
	trinitygo.RegisterController("/v1/metros", metroControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetMetroByID"),
		application.NewRequestMapping(httputil.GET, "", "GetMetroList"),
		application.NewRequestMapping(httputil.POST, "", "CreateMetro"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateMetroByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteMetroByID"),
		application.NewRequestMapping(httputil.DELETE, "", "MultiDeleteMetroByID"),
	)
}

// MetroController metro controller
type MetroController interface {
	GetMetroByID()
	GetMetroList()
	CreateMetro()
	UpdateMetroByID()
	DeleteMetroByID()
	MultiDeleteMetroByID()
}

type metroControllerImpl struct {
	MetroSrv service.MetroService `autowired:"true" resource:"MetroService"`
	Tctx     application.Context  `autowired:"true" transaction:"true"`
}

// GetMetroByID Method
// @Summary Get Metro By ID
// @Description function for MetroController  to get Metro By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/metros/{id} [get]
func (c *metroControllerImpl) GetMetroByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.MetroSrv.GetMetroByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// GetMetroList Method
// @Summary Get Metro list By filter
// @Description function for MetroController  to get Metro list By filter
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/metros [get]
func (c *metroControllerImpl) GetMetroList() {
	res, err := c.MetroSrv.GetMetroList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// CreateMetro Method
// @Summary Create Metro
// @Description function for MetroController  to create Metro
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/metros [post]
func (c *metroControllerImpl) CreateMetro() {
	var newMetro model.Metro
	if err := c.Tctx.GinCtx().BindJSON(&newMetro); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.MetroSrv.CreateMetro(&newMetro)
	c.Tctx.httpResponseCreated(res, err)
	return
}

// UpdateMetroByID Method
// @Summary Modify Metro
// @Description function for MetroController  to Modify Metro
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/metros/{id} [patch]
func (c *metroControllerImpl) UpdateMetroByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.MetroSrv.UpdateMetroByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// DeleteMetroByID Method
// @Summary Delete Metro
// @Description function for MetroController  to delete Metro
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/metros/{id} [delete]
func (c *metroControllerImpl) DeleteMetroByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.MetroSrv.DeleteMetroByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// MultiDeleteMetroByID Method
// @Summary MultiDelete Metro
// @Description function for MetroController  to MultiDelete Metro
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/metros [delete]
func (c *metroControllerImpl) MultiDeleteMetroByID() {
	var deleteBody []modelutil.DeleteParam
	if err := c.Tctx.GinCtx().BindJSON(&deleteBody); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	err := c.MetroSrv.MultiDeleteMetroByID(deleteBody)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
