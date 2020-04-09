package main

import (
	helloworldpb "github.com/PolarPanda611/trinitygo/example/pb/helloworld"

	"github.com/PolarPanda611/trinitygo/example/server/domain/controller/grpc"

	"github.com/PolarPanda611/trinitygo"
	truntime "github.com/PolarPanda611/trinitygo/runtime"
)

func main() {
	trinitygo.SetConfigPath("/Users/daniel/Documents/workspace/trinitygo/example/config/example.toml")
	t := trinitygo.DefaultGRPC()
	t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", true,  func() string { return "" }))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", true,  func() string { return "" }))
	t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", true,  func() string { return "" }))
	t.InitGRPC()
	{
		helloworldpb.RegisterGreeterServer(t.GetGRPCServer(), &grpc.Server{})
	}
	t.ServeGRPC()
}
