package main

import (
	"fmt"
	"os"
	"path"

	"github.com/PolarPanda611/trinitygo"
	"github.com/PolarPanda611/trinitygo/application"
	_ "github.com/PolarPanda611/trinitygo/example/http/domain/controller/http" // init controller
	"github.com/PolarPanda611/trinitygo/example/http/infra"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

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
	currentPath, _ := os.Getwd()
	projectRootPath := path.Join(currentPath, "../")
	configPath := fmt.Sprintf(projectRootPath + "/config/example.toml")
	casbinPath := fmt.Sprintf(projectRootPath + "/config/rbac_with_domains_model.conf")
	trinitygo.SetConfigPath(configPath)
	trinitygo.SetCasbinConfPath(casbinPath)
	trinitygo.SetFuncGetWhoAmI(getUser)
	trinitygo.EnableHealthCheckURL()
	t := trinitygo.DefaultHTTP()
	t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", false, func() string { return uuid.New().String() }, true))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", false, func() string { return "" }, false))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", false, func() string { return "" }, true))
	t.InitHTTP()
	infra.DB = t.DB()
	infra.Migrate()
	t.ServeHTTP()
}

func getUser(app application.Application, c *gin.Context, db *gorm.DB) (interface{}, error) {
	return "dtan11", nil
}
