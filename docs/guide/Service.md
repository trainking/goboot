# Service

- [Service](#service)
	- [1. 概述](#1-概述)
	- [2. 详细设计](#2-详细设计)
		- [2.1 依赖组件](#21-依赖组件)
		- [2.2 Server](#22-server)
		- [2.3 配置文件](#23-配置文件)
		- [2.4 Service](#24-service)
	- [3. 注意事项](#3-注意事项)


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

### 2.3 配置文件

示例：

```yaml
UserService:
  Prefix: "/services"
  Authentication:
    Mode: "TLS"
    CertFile: "../ssldata/cert.crt"
    KeyFile: "../ssldata/private.key"
  Etcd:
    - "192.168.1.9:2379"
```

与服务相关的配置，都是在相关`XXXService`之下:

| 配置项                  | 类型                | 是否必须 | 说明                     |
| :---------------------- | :------------------ | :------- | :----------------------- |
| Prefix                  | string              | 是       | 服务注册的前缀           |
| Authentication          | map[string][string] | 否       | 验证选项                 |
| Authentication.Mode     | string              | 否       | 验证使用方式，可选 `TLS` |
| Authentication.CertFile | string              | 否       | cert文件路径             |
| Authentication.KeyFile  | string              | 否       | 密钥文件的路径           |
| Etcd                    | []string            | 否       | 服务注册的Etcd地址       |


### 2.4 Service

在`pb/proto`中，定义一个RPC:

```protobuf
service UserService {
  rpc GetUserInfo(GetUserInfoArgs) returns (GetUserInfoReply) {}
}

message GetUserInfoArgs {
  int64 user_id = 1; // 用户ID
}
message GetUserInfoReply {
  string user_name = 1; // 用户名
}
```

> 推荐请求message都以Args结尾，返回message都以Reply结尾

使用下面命令，生成Service：

```
protoc --go_out=../ --go_opt=paths=source_relative --go-grpc_out=../ --go-grpc_opt=paths=source_relative *.service.proto
```

也可以使用`booctl`工具生成Service，详情参考[bootctl](./bootctl.md)。

## 3. 注意事项