package db

import (
	"github.com/PolarPanda611/trinitygo/logger"

	"github.com/PolarPanda611/trinitygo/application"
)

type dbLogger struct {
	app    application.Application
	config *logger.Config
}
