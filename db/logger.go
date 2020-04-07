package db

import (
	"trinitygo/application"
	"trinitygo/logger"
)

type dbLogger struct {
	app    application.Application
	config *logger.Config
}
