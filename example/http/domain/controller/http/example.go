package http

import (
	"errors"
	"strconv"
	"sync"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/example/http/domain/service"
	"github.com/gin-gonic/gin"
)

func init() {
	trinitygo.BindController("GET@/ping/:id/*dd", &sync.Pool{
		New: func() interface{} {
			service := new(Server)
			return service
		},
	})
}

// Server  test
type Server struct {
	Service service.Service
}

// Get Method
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
