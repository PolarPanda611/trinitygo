package http

import (
	"errors"
	"strconv"
	"sync"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/example/http/domain/service"
	"github.com/PolarPanda611/trinitygo/httputils"
	"github.com/gin-gonic/gin"
)

func init() {
	trinitygo.BindController("/ping",
		&sync.Pool{
			New: func() interface{} {
				service := new(Server)
				return service
			},
		},
		httputils.NewRequestMapping(httputils.GET, "/:id", "Get"),
		httputils.NewRequestMapping(httputils.GET, "", "Getsssss"),
	)

}

// Server  test
type Server struct {
	Service service.Service
	C       *gin.Context
	Tctx    application.Context
}

// @Summary get user by id
// @Description get user by id test
// @Produce  json
// @Param   id     path    int64     true        "user id"
// @Success 200 {string} json "{"Status":200,"Result":{},"Runtime":"ok"}"
// @Failure 400 {string} json "{"Status":400,"Result":{},"Runtime":"ok"}"
// @Router /ping/{id} [get]
func (s *Server) Get() {
	id, err := strconv.Atoi(s.C.Param("id"))
	if err != nil {
		s.Tctx.Response(400, nil, errors.New("wrong id"))
		return
	}
	res, err := s.Service.Get(id)
	if err != nil {
		s.Tctx.Response(400, nil, err)
		return
	}
	s.Tctx.Response(200, res, nil)
	return

}

// Getsssss Method
func (s *Server) Getsssss() {

	res, err := s.Service.Get(200)
	if err != nil {
		s.Tctx.Response(400, nil, err)
		return
	}
	s.Tctx.Response(200, res, nil)
	return
}
