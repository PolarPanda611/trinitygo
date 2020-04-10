package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/PolarPanda611/trinitygo"
	_ "github.com/PolarPanda611/trinitygo/example/http/domain/controller/http" // init controller
	truntime "github.com/PolarPanda611/trinitygo/runtime"
	"github.com/google/uuid"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	projectRootPath := path.Join(path.Dir(b), "../")
	configPath := fmt.Sprintf(projectRootPath + "/config/example.toml")
	trinitygo.SetConfigPath(configPath)
	t := trinitygo.DefaultHTTP()
	t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", false, func() string { return uuid.New().String() }))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", false, func() string { return "124" }))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", false, func() string { return "dtan11" }))
	t.InitHTTP()
	t.ServeHTTP()
}
