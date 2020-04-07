package main

import (
	helloworldpb "github.com/PolarPanda611/trinitygo/example/pb/helloworld"

	"github.com/PolarPanda611/trinitygo/example/server/domain/controller/grpc"

	"github.com/PolarPanda611/trinitygo"
)

func main() {
	trinitygo.SetConfigPath("/Users/daniel/Documents/workspace/trinitygo/example/config/example.toml")
	t := trinitygo.DefaultGRPC()
	t.InitGRPC()
	{
		helloworldpb.RegisterGreeterServer(t.GetGRPCServer(), &grpc.Server{})
	}
	t.ServeGRPC()
}
