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

type userAuthorization struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var _ UserController = new(userControllerImpl)

func init() {
	trinitygo.RegisterController("/v1/users", userControllerImpl{},
		application.NewRequestMapping(httputil.GET, "/:id", "GetUserByID"),
		application.NewRequestMapping(httputil.GET, "", "GetUserList"),
		application.NewRequestMapping(httputil.POST, "", "CreateUser"),
		application.NewRequestMapping(httputil.PATCH, "/:id", "UpdateUserByID"),
		application.NewRequestMapping(httputil.DELETE, "/:id", "DeleteUserByID"),
		application.NewRequestMapping(httputil.DELETE, "", "MultiDeleteUserByID"),
	)
	trinitygo.RegisterController("/v1/currentUser", userControllerImpl{},
		application.NewRequestMapping(httputil.GET, "", "GetCurrentUser"),
	)
	trinitygo.RegisterController("/v1/login", userControllerImpl{},
		application.NewRequestMapping(httputil.POST, "", "VerifyUserByAuthorization"),
	)

}

// UserController user controller
type UserController interface {
	GetCurrentUser()
	GetUserByID()
	GetUserList()
	CreateUser()
	UpdateUserByID()
	DeleteUserByID()
	MultiDeleteUserByID()
}

type userControllerImpl struct {
	UserSrv service.UserService `autowired:"true" resource:"UserService"`
	Tctx    application.Context `autowired:"true" transaction:"true"`
}

func (c *userControllerImpl) VerifyUserByAuthorization() {
	var userAuth userAuthorization
	if err := c.Tctx.GinCtx().BindJSON(&userAuth); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	_, err := c.UserSrv.GetUserByUserName(userAuth.UserName)
	if err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6Ik1BSU4iLCJwaS5hdG0iOiI2In0.eyJzY29wZSI6WyJvcGVuaWQiLCJwcm9maWxlIl0sImNsaWVudF9pZCI6IkNjZjk4ODBjZThiYThlYmUxMGM3ZDY3YjgyNTEzYjQxM2E2MWFkMDI5IiwiaXNzIjoiaWRwZGVjYXRobG9uLnByZXByb2Qub3JnIiwianRpIjoieG1FMjVNNVVnYyIsInN1YiI6IkRUQU4xMSIsInVpZCI6IkRUQU4xMSIsIm9yaWdpbiI6ImNvcnBvcmF0ZSIsImV4cCI6MTU4NTQ5MTAyOH0.C1WZLwbQIOD1FEoHZvz2x3Je1_m2b-rcu4IYAYR-3x64bS5tXasS9z9Qf3kDafUKqqYZDcb4kpSoirFzdzvVWUo4zeOUYv1OYkxpFqxQia4YrHZGPC5VQi9lVaYNEm4CBCYV3ZoM-sJh1spcpY0x5g-Z9nZhJCsd9l6m6JQuF0aRh2zlEir04qKgk67m25zeGA1s8tAm27hEOkhtS_00-MlFrsMNO7ZS50xdT-gQ0_20_VVYMVhbZAsQIi3mxTYag-bVLrQFLTLLIbNW1nFuH0P6xOYdSW6WuBsI23HRJ5HeTIOcdOSKWi-6MtvEuN4vS69k_zg-oj2rsTlQCPoT2Q"
	res := map[string]interface{}{
		// "data":  user,
		"token": token,
	}
	c.Tctx.HTTPResponseOk(res, err)
	return

}
func (c *userControllerImpl) GetCurrentUser() {
	id, _ := strconv.ParseInt(c.Tctx.Runtime()["user_id"], 10, 64)
	res, err := c.UserSrv.GetUserByID(id)
	c.Tctx.HTTPResponseOk(res, err)
	return
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
func (c *userControllerImpl) GetUserByID() {
	id, _ := strconv.ParseInt(c.Tctx.GinCtx().Params.ByName("id"), 10, 64)
	res, err := c.UserSrv.GetUserByID(id)
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
func (c *userControllerImpl) GetUserList() {
	res, err := c.UserSrv.GetUserList(c.Tctx.GinCtx().Request.URL.RawQuery)
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
func (c *userControllerImpl) CreateUser() {
	var newUser model.User
	if err := c.Tctx.GinCtx().BindJSON(&newUser); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	res, err := c.UserSrv.CreateUser(&newUser)
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

// MultiDeleteUserByID Method
// @Summary multi delete user by id
// @Description function for controller to multi delete user list
// @accept  json
// @Produce json
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /metro_api/v1/users [delete]
func (c *userControllerImpl) MultiDeleteUserByID() {
	var deleteBody []modelutil.DeleteParam
	if err := c.Tctx.GinCtx().BindJSON(&deleteBody); err != nil {
		c.Tctx.HTTPResponseInternalErr(err)
		return
	}
	err := c.UserSrv.MultiDeleteUserByID(deleteBody)
	c.Tctx.HTTPResponseOk(nil, err)
	return
}
