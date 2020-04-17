package http

import (
	"errors"
	"strconv"
	"sync"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/example/http/domain/service"
	"github.com/PolarPanda611/trinitygo/httputil"
	"github.com/PolarPanda611/trinitygo/util"
)

var _ UserController = new(userControllerImpl)

func init() {
	trinitygo.BindController("/users",
		&sync.Pool{
			New: func() interface{} {
				controller := new(userControllerImpl)
				return controller
			},
		},
		application.NewRequestMapping(httputil.GET, "/:id", "GET", PermissionValidator([]string{"manager"}), gValidator, g1Validator),
		// application.NewRequestMapping(httputil.GET, "/:id", "GET"),
		application.NewRequestMapping(httputil.GET, "", "Getsssss"),
	)

}

// UserController user controller
type UserController interface {
	GET()
	Getsssss()
}

// UserController  test
type userControllerImpl struct {
	UserSrv service.UserService
	Tctx    application.Context
}

// GET get func
// @Summary get user by id
// @Description get user by id test
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "user id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Router /users/{id} [get]
func (s *userControllerImpl) GET() {
	id, _ := strconv.Atoi(s.Tctx.GinCtx().Param("id"))
	res, err := s.UserSrv.GetUserByID(id)
	s.Tctx.HTTPResponseOk(res, err)
	return
}

var gValidator = func(tctx application.Context) {
	id, _ := strconv.Atoi(tctx.GinCtx().Param("id"))
	if id < 3 {
		tctx.HTTPResponseUnauthorizedErr(errors.New("gValidator no permission"))
	}
	return
}

var g1Validator = func(tctx application.Context) {
	id, _ := strconv.Atoi(tctx.GinCtx().Param("id"))
	if id > 3 {
		tctx.HTTPResponseUnauthorizedErr(errors.New("g1Validator no permission"))
	}
	return
}

// PermissionValidator example validator
func PermissionValidator(requiredP []string) func(application.Context) {
	return func(c application.Context) {
		// c.GinCtx().Set("permission", []string{"employee"}) // no permission
		c.GinCtx().Set("permission", []string{"employee", "manager"}) // ok
		in := util.SliceInSlice(requiredP, c.GinCtx().GetStringSlice("permission"))
		if !in {
			c.HTTPResponseUnauthorizedErr(errors.New("np permission"))
		}
	}
}

// Getsssss Method
// @Summary get user list
// @Description get user by id test
// @accept  json
// @Produce json
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Router /users [get]
func (s *userControllerImpl) Getsssss() {
	res, err := s.UserSrv.GetUserListByQuery(s.Tctx.GinCtx().Request.URL.RawQuery)
	s.Tctx.HTTPResponseOk(res, err)
	return
}
