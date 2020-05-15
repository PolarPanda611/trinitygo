package main

import (
	"fmt"
	"os"

	_ "test/domain/controller/http"

	"github.com/PolarPanda611/trinitygo"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/google/uuid"
)

// @title test
// @version 1.0
// @description  test
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
	trinitygo.EnableHealthCheckURL()
	t := trinitygo.DefaultHTTP()
	t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", false, func() string { return uuid.New().String() }, true))
	t.InitHTTP()
	t.ServeHTTP()
}
