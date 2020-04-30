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

var _ UserController = new(userControllerImpl)

func init() {
	trinitygo.RegisterController("/users", userControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetUserByID"),
		application.NewRequestMapping(httputil.GET, "", "GetUserList"),
		application.NewRequestMapping(httputil.POST, "", "CreateUser"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateUserByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteUserByID"),
	)
}

// UserController user controller
type UserController interface {
	GetUserByID()
	GetUserList()
	CreateUser()
	UpdateUserByID()
	DeleteUserByID()
}

type userControllerImpl struct {
	UserSrv service.UserService `autowired:"true" resource:"UserService"`
	Tctx    application.Context `autowired:"true" transaction:"true"`
}

func (c *userControllerImpl) GetUserByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.UserSrv.GetUserByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

func (c *userControllerImpl) GetUserList() {
	res, err := c.UserSrv.GetUserList(c.Tctx.GinCtx().Request.URL.RawQuery)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

func (c *userControllerImpl) CreateUser() {
	var newUser model.User
	if err := c.Tctx.GinCtx().BindJSON(&newUser); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.UserSrv.CreateUser(&newUser)
	c.Tctx.HTTPResponseOk(res, err)
	return
}

func (c *userControllerImpl) UpdateUserByID() {
	change, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.UserSrv.UpdateUserByID(id, dVersion, change)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}

func (c *userControllerImpl) DeleteUserByID() {
	_, dVersion, err := util.DecodeReqBodyToMap(c.Tctx.GinCtx())
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	err = c.UserSrv.DeleteUserByID(id, dVersion)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
