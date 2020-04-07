package main

import (
	"trinitygo"
	helloworldpb "trinitygo/example/pb/helloworld"
	"trinitygo/example/server/domain/controller/grpc"
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
