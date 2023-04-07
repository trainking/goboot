# Service

## 1. 概述

**goboot**中的Service是通过`gRPC`定制开发的微服务实现，主要实现都在`boot/service`包中，提供了一个基础的`BaseService`实现。因为使用`Etcd`服务注册和发现，所以需要依赖ETCD。

## 2. 详细设计

### 2.1 依赖组件

**goboot**开发微服务，需要以下几个组件：

* etcd：服务注册和服务发现中间件
* protoc：编译protobuf文件工具
* proto-gen-go：protobuf的go编译插件
* proto-gen-go-grpc：protobuf的gRPC的编译插件

### 2.2 Server

**goboot**中，每一个微服务都需要实现一个Server结构体，继承（组合）了`BaseService`，实现对应定义的gRPC接口：

```go

type (
	Server struct {
		pb.UnimplementedUserServiceServer
		boot.BaseService
	}
)

// New 创建Server， 使用配置文件
func New(name string, configPath string, addr string, instanceId int64) boot.Instance {
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	server := new(Server)
	server.Name = name
	server.Config = v
	server.Addr = addr
	server.IntanceID = instanceId
	return server
}

// Init 服务初始化
func (s *Server) Init() error {
	var err error
	if err = s.BaseService.Init(); err != nil {
		log.Errorf("BaseServcie error %v", err)
		return err
	}

	pb.RegisterUserServiceServer(s.GrpcServer, s)

	return nil
}

```

Server是Instance的实现，所以它必须要有以下参数：

| 参数       | Options               | 类型   | 说明                                                                          |
| :--------- | :-------------------- | :----- | :---------------------------------------------------------------------------- |
| name       | -name,--name          | string | 名称，定义此实例作用，可以根据name组合Module实现不同作用的实例，如Gate, Lobby |
| addr       | -addr, --addr         | string | 监听的ip:port，如`127.0.0.1:6001`                                             |
| configPath | -config, --config     | string | 配置文件的路径                                                                |
| instanceId | -instance, --instance | int64  | 实例的ID，实例的唯一标识                                                      |


## 3. 注意事项