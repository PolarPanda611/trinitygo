package application

import (
	"fmt"
	"time"

	"github.com/PolarPanda611/trinitygo/logger"
)

// Logger db logger
type Logger struct {
	app     Application
	config  *logger.Config
	runtime map[string]string
}

// NewDBLogger new db logger
func NewDBLogger(app Application, runtime map[string]string, config ...*logger.Config) *Logger {
	c := logger.DefaultConfig()
	if len(config) > 0 {
		c = config[0]
	}
	return &Logger{
		app:     app,
		config:  c,
		runtime: runtime,
	}

}

// Print db logger func
func (l *Logger) Print(v ...interface{}) {
	dblogLevel, _ := v[0].(string)
	if dblogLevel == "sql" {
		// logInterface = append(logInterface, "DBRunningFile=")
		// logInterface = append(logInterface, fmt.Sprint(v[1]))
		// logInterface = append(logInterface, "DBRunningTime=")
		//
		// logInterface = append(logInterface, DBRunningTime)
		// logInterface = append(logInterface, "DBSQL=")
		// logInterface = append(logInterface, fmt.Sprint(v[3]))
		// logInterface = append(logInterface, "DBParams=")
		// logInterface = append(logInterface, fmt.Sprint(v[4]))
		// logInterface = append(logInterface, "DBEffectedRows=")
		// logInterface = append(logInterface, fmt.Sprint(v[5]))

		line := ""
		if l.config.Runtime {
			for _, v := range l.app.RuntimeKeys() {
				line += fmt.Sprintf("%v %v ", v.GetKeyName(), l.runtime[v.GetKeyName()])
			}
		}
		line += fmt.Sprintf("%v %v ", "DBRunningFile", v[1])
		DBRunningTime, _ := v[2].(time.Duration)
		line += fmt.Sprintf("%v %v ", "DBRunningTime", DBRunningTime)
		line += fmt.Sprintf("%v %v ", "DBSQL", v[3])
		line += fmt.Sprintf("%v %v ", "DBParams", v[4])
		line += fmt.Sprintf("%v %v ", "DBEffectedRows", v[5])

		l.app.Logger().Info(line)
	}
}
