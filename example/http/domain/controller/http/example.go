package http

import (
	"errors"
	"strconv"
	"sync"

	"github.com/PolarPanda611/trinitygo"
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
		httputils.NewRequestMapping(httputils.GET, "/:id", "GET"),
		httputils.NewRequestMapping(httputils.GET, "", "Getsssss"),
	)

}

// Server  test
type Server struct {
	Service service.Service
}

// GET Method
func (s *Server) GET(c *gin.Context) (int, interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 400, nil, errors.New("wrong id")
	}
	res, err := s.Service.Get(id)
	if err != nil {
		return 400, nil, err
	}
	return 200, res, nil

}

// Getsssss Method
func (s *Server) Getsssss(c *gin.Context) (int, interface{}, error) {

	res, err := s.Service.Get(200)
	if err != nil {
		return 400, nil, err
	}
	return 200, res, nil
}
