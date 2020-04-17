golang restframework plugin with gin+gorm, fast and high scalable    
p.s: spring restframework like :)

## 安装

打开终端输入

```bash
$ go get -u github.com/PolarPanda611/trinitygo
```

done.

## 特性

* 集成gorm
* 集成gin
* 快速注册路由
* 链路追踪，返回错误代码行数,及sql
* 支持请求事务
* 支持自定义用户权限查询
* 支持自定义接口访问权限
* 支持自定义接口数据访问权限
* 自定义搜索
* 自定义预加载（gorm preload）
* 自定义排序
* 自定义查询包含字段

## 文档

## Http Server 
```new HTTP server 
// can see the example in example/http

trinitygo.SetConfigPath(configPath)  // put config path here
t := trinitygo.DefaultHTTP()
t.InitHTTP()
t.ServeHTTP()

```






## GRPC Server 
```new GRPC server 
// can see the example in example/server

trinitygo.SetConfigPath(configPath) // put config path here
t := trinitygo.DefaultGRPC()
t.RegRuntimeKey(truntime.NewRuntimeKey("trace_id", true, func() string { return "" }, true))
t.RegRuntimeKey(truntime.NewRuntimeKey("user_id", true, func() string { return "" }, true))
t.RegRuntimeKey(truntime.NewRuntimeKey("user_name", true, func() string { return "" }, true))
t.InitGRPC()
{
    helloworldpb.RegisterGreeterServer(t.GetGRPCServer(), &grpc.Server{})  // register your grpc server here
}
t.ServeGRPC()

```


