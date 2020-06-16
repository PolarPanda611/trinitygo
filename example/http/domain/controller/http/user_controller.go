package http

import (
	"http/domain/model"
	"http/domain/service"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	modelutil "github.com/PolarPanda611/trinitygo/crud/model"
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
		application.NewRequestMapping(httputil.DELETE, "", "MultiDeleteUserByID"),
	)
	trinitygo.RegisterController("/benchmark", userControllerImpl{},
		application.NewRequestMapping(httputil.GET, "", "Benchmark"),
	)
}

// UserController user controller
type UserController interface {
	GetUserByID(args struct {
		ID int64 `path_param:"id"`
	}) (*model.User, error)
	GetUserList(args struct {
		Query string `query_param:""`
	}) (interface{}, error)
	CreateUser(args struct {
		User model.User `body_param:""`
	}) (*model.User, error)
	UpdateUserByID(args struct {
		ID       int64                  `path_param:"id"`
		Change   map[string]interface{} `body_param:""`
		DVersion string                 `body_param:"d_version"`
	}) error
	DeleteUserByID(args struct {
		ID       int64  `path_param:"id"`
		DVersion string `body_param:"d_version"`
	}) error
	MultiDeleteUserByID(args struct {
		DeleteParamList []modelutil.DeleteParam `body_param:""`
	}) error
	Benchmark() error
}

type userControllerImpl struct {
	UserSrv service.UserService `autowired:"true" resource:"UserService"`
	Tctx    application.Context `autowired:"true" transaction:"true"`
}

// GetUserByID Method
// @Summary Get User By ID
// @Description function for UserController  to get User By ID
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /http/users/{id} [get]
func (c *userControllerImpl) GetUserByID(args struct {
	ID int64 `path_param:"id"`
}) (*model.User, error) {
	return c.UserSrv.GetUserByID(args.ID)
}

// GetUserList Method
// @Summary Get User list By filter
// @Description function for UserController  to get User list By filter
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /http/users [get]
func (c *userControllerImpl) GetUserList(args struct {
	Query string `query_param:""`
}) (interface{}, error) {
	return c.UserSrv.GetUserList(args.Query)
}

// CreateUser Method
// @Summary Create User
// @Description function for UserController  to create User
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 201 {string} json "{"Status":201,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /http/users [post]
func (c *userControllerImpl) CreateUser(args struct {
	User model.User `body_param:""`
}) (*model.User, error) {
	c.Tctx.HTTPStatus(201)
	return c.UserSrv.CreateUser(&args.User)
}

// UpdateUserByID Method
// @Summary Modify User
// @Description function for UserController  to Modify User
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /http/users/{id} [patch]
func (c *userControllerImpl) UpdateUserByID(args struct {
	ID       int64                  `path_param:"id"`
	Change   map[string]interface{} `body_param:""`
	DVersion string                 `body_param:"d_version"`
}) error {
	return c.UserSrv.UpdateUserByID(args.ID, args.DVersion, args.Change)
}

// DeleteUserByID Method
// @Summary Delete User
// @Description function for UserController  to delete User
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "id"
// @Param  q  query string false "name search by q" Format(email)
// @Success 204 {string} json "{"Status":204,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /http/users/{id} [delete]
func (c *userControllerImpl) DeleteUserByID(args struct {
	ID       int64  `path_param:"id"`
	DVersion string `body_param:"d_version"`
}) error {
	c.Tctx.HTTPStatus(204)
	return c.UserSrv.DeleteUserByID(args.ID, args.DVersion)
}

// MultiDeleteUserByID Method
// @Summary MultiDelete User
// @Description function for UserController  to MultiDelete User
// @accept  json
// @Produce json
// @Param  q  query string false "name search by q" Format(email)
// @Success 204 {string} json "{"Status":204,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Security ApiKeyAuth
// @Router /http/users [delete]
func (c *userControllerImpl) MultiDeleteUserByID(args struct {
	DeleteParamList []modelutil.DeleteParam `body_param:""`
}) error {
	c.Tctx.HTTPStatus(204)
	return c.UserSrv.MultiDeleteUserByID(args.DeleteParamList)
}

func (c *userControllerImpl) Benchmark() error {
	return nil
}
