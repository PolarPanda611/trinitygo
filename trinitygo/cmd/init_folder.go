package cmd

import "os"

func initHTTPFolder(rootPath string) {
	os.MkdirAll(rootPath, os.ModePerm)
	os.MkdirAll(rootPath+"/conf", os.ModePerm)
	os.MkdirAll(rootPath+"/domain", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/entity", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/object", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/repository", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/service", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/controller", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/controller/http", os.ModePerm)

}

func initGRPCFolder(rootPath string) {
	os.MkdirAll(rootPath, os.ModePerm)
	os.MkdirAll(rootPath+"/conf", os.ModePerm)
	os.MkdirAll(rootPath+"/domain", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/entity", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/object", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/repository", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/service", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/controller", os.ModePerm)
	os.MkdirAll(rootPath+"/domain/controller/grpc", os.ModePerm)

}
