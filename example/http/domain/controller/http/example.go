package http

import (
	"errors"
	"strconv"
	"sync"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/example/http/domain/service"
	"github.com/PolarPanda611/trinitygo/httputils"
	"github.com/PolarPanda611/trinitygo/utils"
)

func init() {
	trinitygo.BindController("/ping",
		&sync.Pool{
			New: func() interface{} {
				service := new(Server)
				return service
			},
		},
		application.NewRequestMapping(httputils.GET, "/:id", "", PermissionValidator([]string{"manager"}), gValidator),
		application.NewRequestMapping(httputils.GET, "", "Getsssss"),
	)

}

// Server  test
type Server struct {
	Service service.Service
	Tctx    application.Context
}

// Get get func
// @Summary get user by id
// @Description get user by id test
// @accept  json
// @Produce json
// @Param   id     path    int64     true        "user id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Router /ping/{id} [get]

func (s *Server) GET() {
	id, _ := strconv.Atoi(s.Tctx.GinCtx().Param("id"))
	res, err := s.Service.Get(id)
	s.Tctx.HTTPResponseOk(res, err)
	return
}

var gValidator = func(tctx application.Context) {
	id, _ := strconv.Atoi(tctx.GinCtx().Param("id"))
	if id < 100 {
		tctx.HTTPResponseUnauthorizedErr(errors.New("gValidator no permission"))
	}
	return
}

var g1Validator = func(tctx application.Context) {
	id, _ := strconv.Atoi(tctx.GinCtx().Param("id"))
	if id > 200 {
		tctx.HTTPResponseUnauthorizedErr(errors.New("g1Validator no permission"))
	}
	return
}

// PermissionValidator example validator
func PermissionValidator(requiredP []string) func(application.Context) {
	return func(c application.Context) {
		// c.GinCtx().Set("permission", []string{"employee"}) // no permission
		c.GinCtx().Set("permission", []string{"employee", "manager"}) // ok
		in := utils.SliceInSlice(requiredP, c.GinCtx().GetStringSlice("permission"))
		if !in {
			c.HTTPResponseUnauthorizedErr(errors.New("np permission"))
		}
	}
}

// Getsssss Method
func (s *Server) Getsssss() {

	res, err := s.Service.Get(200)
	s.Tctx.HTTPResponseOk(res, err)
	return
}
