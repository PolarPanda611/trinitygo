package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/PolarPanda611/trinitygo"
	_ "github.com/PolarPanda611/trinitygo/example/http/domain/controller/http" // init controller
	"github.com/PolarPanda611/trinitygo/example/http/infra"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/google/uuid"

	_ "github.com/PolarPanda611/trinitygo/example/http/docs"
)

// @title Trinity HTTP Example API
// @version 1.0
// @description This is a sample trinity http server
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8088
// @BasePath /trinitygo/
func main() {
	_, b, _, _ := runtime.Caller(0)
	projectRootPath := path.Join(path.Dir(b), "../")
	configPath := fmt.Sprintf(projectRootPath + "/config/example.toml")
	trinitygo.SetConfigPath(configPath)
	t := trinitygo.DefaultHTTP()
	t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", false, func() string { return uuid.New().String() }, true))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", false, func() string { return "" }, false))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", false, func() string { return "" }, true))
	t.InitHTTP()
	infra.DB = t.DB()
	infra.Migrate()
	t.ServeHTTP()
}
