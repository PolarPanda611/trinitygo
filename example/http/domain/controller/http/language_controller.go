package http

import (
	"strconv"

	"github.com/PolarPanda611/trinitygo/example/http/domain/model"

	"github.com/PolarPanda611/trinitygo/example/http/domain/service"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/crud/util"
	"github.com/PolarPanda611/trinitygo/httputil"
)

var _ LanguageController = new(languageControllerImpl)

func init() {
	trinitygo.RegisterController("/languages", languageControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetLanguageByID"),
		application.NewRequestMapping(httputil.GET, "", "GetLanguageList"),
		application.NewRequestMapping(httputil.POST, "", "CreateLanguage"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateLanguageByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteLanguageByID"),
	)
}

// LanguageController language controller
type LanguageController interface {
	GetLanguageByID()
	GetLanguageList()
	CreateLanguage()
	UpdateLanguageByID()
	DeleteLanguageByID()
}

type languageControllerImpl struct {
	LanguageSrv service.LanguageService `autowired:"true" resource:"LanguageService"`
	Tctx        application.Context     `autowired:"true" transaction:"true"`
}

func (c *languageControllerImpl) GetLanguageByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.LanguageSrv.GetLanguageByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

func (c *languageControllerImpl) GetLanguageList() {
	res, err := c.LanguageSrv.GetLanguageList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

func (c *languageControllerImpl) CreateLanguage() {
	var newLanguage model.Language
	if err := c.Tctx.GinCtx().BindJSON(&newLanguage); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.LanguageSrv.CreateLanguage(&newLanguage)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

func (c *languageControllerImpl) UpdateLanguageByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.LanguageSrv.UpdateLanguageByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

func (c *languageControllerImpl) DeleteLanguageByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.LanguageSrv.DeleteLanguageByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
