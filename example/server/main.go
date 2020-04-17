package main

import (
	"fmt"
	"os"
	"path"

	helloworldpb "github.com/PolarPanda611/trinitygo/example/pb/helloworld"

	"github.com/PolarPanda611/trinitygo/example/server/domain/controller/grpc"

	"github.com/PolarPanda611/trinitygo"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
)

func main() {
	currentPath, _ := os.Getwd()
	projectRootPath := path.Join(currentPath, "../")
	configPath := fmt.Sprintf(projectRootPath + "/config/example.toml")
	trinitygo.SetConfigPath(configPath)
	t := trinitygo.DefaultGRPC()
	t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", true, func() string { return "" }, true))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", true, func() string { return "" }, true))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", true, func() string { return "" }, true))
	t.InitGRPC()
	{
		helloworldpb.RegisterGreeterServer(t.GetGRPCServer(), &grpc.Server{})
	}
	t.ServeGRPC()
}
