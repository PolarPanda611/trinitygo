package main

import (
	"fmt"
	_ "http/docs"
	_ "http/domain/controller/http"
	"http/infra/db"
	"http/infra/migrate"
	"os"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/keyword"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/google/uuid"
)

// @title http
// @version 1.0
// @description  http
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host address:port
// @BasePath /trinitygo/
func main() {
	currentPath, _ := os.Getwd()
	configPath := fmt.Sprintf(currentPath + "/conf/conf.toml")
	trinitygo.SetConfigPath(configPath)
	trinitygo.SetResponseFactory(CustomizeResponseFactory)
	trinitygo.SetKeyword(keyword.Keyword{
		SearchBy:      "SearchBy",
		PageNum:       "current",
		PageSize:      "pageSize",
		OrderBy:       "OrderBy",
		PaginationOff: "PaginationOff",
	})
	trinitygo.EnableHealthCheckURL()
	t := trinitygo.DefaultHTTP()
	t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", false, func() string { return uuid.New().String() }, true))
	t.InitHTTP()
	db.DB = t.DB()
	migrate.Migrate()
	t.ServeHTTP()
}

// CustomizeResponseFactory customize response formatter
func CustomizeResponseFactory(status int, res interface{}, runtime map[string]string) interface{} {
	resMap, ok := res.(map[string]interface{})
	if !ok {
		resMap = make(map[string]interface{})
		resMap["data"] = res
	}
	resMap["status"] = status
	for k, v := range runtime {
		resMap[k] = v
	}
	return resMap
}
