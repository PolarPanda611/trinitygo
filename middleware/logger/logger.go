package logger

import (
	"github.com/PolarPanda611/trinitygo/application"
	"github.com/PolarPanda611/trinitygo/logger"
	"github.com/gin-gonic/gin"
)

type middlewareLogger struct {
	app    application.Application
	config *logger.Config
}

// New record log
func New(app application.Application, config ...*logger.Config) gin.HandlerFunc {
	c := logger.DefaultConfig()
	if len(config) > 0 {
		c = config[0]
	}

	l := &middlewareLogger{app: app, config: c}
	return l.LogMiddleware()
}

func (m *middlewareLogger) LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
