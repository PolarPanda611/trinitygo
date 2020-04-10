package logger

import (
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/gin-gonic/gin"
)

// New record log
func New(app application.Application, config ...*application.LogConfig) gin.HandlerFunc {
	c := application.DefaultLogConfig()
	if len(config) > 0 {
		c = config[0]
	}

	l := application.NewLogLogger(app, c)
	return l.Middleware()
}
