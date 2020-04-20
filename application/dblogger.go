package application

import (
	"fmt"
	"time"
)

// DBLogger db logger
type DBLogger struct {
	app     Application
	config  *LogConfig
	runtime map[string]string
}

// NewDBLogger new db logger
func NewDBLogger(app Application, runtime map[string]string, config ...*LogConfig) *DBLogger {
	c := DefaultLogConfig()
	if len(config) > 0 {
		c = config[0]
	}
	return &DBLogger{
		app:     app,
		config:  c,
		runtime: runtime,
	}

}

// Print db logger func
func (l *DBLogger) Print(v ...interface{}) {
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
				if v.IsLog() {
					line += fmt.Sprintf("%v %v ", v.GetKeyName(), l.runtime[v.GetKeyName()])
				}
			}
		}
		DBRunningTime, _ := v[2].(time.Duration)
		line += fmt.Sprintf("%v ", DBRunningTime)
		line += fmt.Sprintf("%v ", v[1])

		line += fmt.Sprintf("%v ", v[3])
		line += fmt.Sprintf("%v ", v[4])
		line += fmt.Sprintf("%v ", v[5])

		l.app.Logger().Info(line)
	}
}
