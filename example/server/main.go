package main

import (
	"trinitygo"
	"trinitygo/db"
	helloworldpb "trinitygo/example/pb/helloworld"
	"trinitygo/example/server/domain/controller/grpc"

	"github.com/jinzhu/gorm"
)

func main() {
	t := trinitygo.DefaultGRPC()
	t.InstallDB(func() *gorm.DB {
		return db.DefaultInstallGORM(
			true,
			true,
			"postgres",
			"asset_",
			"host=127.0.0.1 port=60901 user=trinity password= dbname=trinity sslmode=disable",
			10,
			100,
		)
	})
	t.InitGRPC()
	{
		helloworldpb.RegisterGreeterServer(t.GetGRPCServer(), &grpc.Server{})
	}
	t.ServeGRPC(":50051")
}
