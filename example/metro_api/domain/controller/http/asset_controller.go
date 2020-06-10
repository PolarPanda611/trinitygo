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

var _ AssetController = new(assetControllerImpl)

func init() {
	trinitygo.RegisterController("/assets", assetControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetAssetByID"),
		application.NewRequestMapping(httputil.GET, "", "GetAssetList"),
		application.NewRequestMapping(httputil.POST, "", "CreateAsset"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateAssetByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteAssetByID"),
		application.NewRequestMapping(httputil.DELETE, "", "MultiDeleteAssetByID"),
	)
}

// AssetController asset controller
type AssetController interface {
	GetAssetByID()
	GetAssetList()
	CreateAsset()
	UpdateAssetByID()
	DeleteAssetByID()
	MultiDeleteAssetByID()
}

type assetControllerImpl struct {
	AssetSrv service.AssetService `autowired:"true" resource:"AssetService"`
	Tctx     application.Context  `autowired:"true" transaction:"true"`
}

// GetAssetByID Method
// @Summary Get Asset By ID
// @Description function for AssetController  to get Asset By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/assets/{id} [get]
func (c *assetControllerImpl) GetAssetByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.AssetSrv.GetAssetByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// GetAssetList Method
// @Summary Get Asset list By filter
// @Description function for AssetController  to get Asset list By filter
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/assets [get]
func (c *assetControllerImpl) GetAssetList() {
	res, err := c.AssetSrv.GetAssetList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

// CreateAsset Method
// @Summary Create Asset
// @Description function for AssetController  to create Asset
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/assets [post]
func (c *assetControllerImpl) CreateAsset() {
	var newAsset model.Asset
	if err := c.Tctx.GinCtx().BindJSON(&newAsset); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.AssetSrv.CreateAsset(&newAsset)
	c.Tctx.httpResponseCreated(res, err)
	return
}

// UpdateAssetByID Method
// @Summary Modify Asset
// @Description function for AssetController  to Modify Asset
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/assets/{id} [patch]
func (c *assetControllerImpl) UpdateAssetByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.AssetSrv.UpdateAssetByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// DeleteAssetByID Method
// @Summary Delete Asset
// @Description function for AssetController  to delete Asset
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/assets/{id} [delete]
func (c *assetControllerImpl) DeleteAssetByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.AssetSrv.DeleteAssetByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

// MultiDeleteAssetByID Method
// @Summary MultiDelete Asset
// @Description function for AssetController  to MultiDelete Asset
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/assets [delete]
func (c *assetControllerImpl) MultiDeleteAssetByID() {
	var deleteBody []modelutil.DeleteParam
	if err := c.Tctx.GinCtx().BindJSON(&deleteBody); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	err := c.AssetSrv.MultiDeleteAssetByID(deleteBody)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
