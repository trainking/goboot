# Service的最佳实践

- [Service的最佳实践](#service的最佳实践)
	- [1. 概述](#1-概述)
	- [2. 基础概念](#2-基础概念)
		- [2.1 依赖组件](#21-依赖组件)
		- [2.2 Server](#22-server)
		- [2.3 配置文件](#23-配置文件)
	- [3. 实践](#3-实践)
		- [3.1 建立一个Service](#31-建立一个service)
			- [3.1.1 配置proto](#311-配置proto)
			- [3.1.2 增加入口文件和Server](#312-增加入口文件和server)
			- [3.1.3 创建一个client](#313-创建一个client)
			- [3.1.4 实现接口](#314-实现接口)
		- [3.1 逻辑处理](#31-逻辑处理)
		- [3.2 Model数据处理](#32-model数据处理)
	- [4. 注意事项](#4-注意事项)


## 1. 概述

**goboot**中的Service是通过`gRPC`定制开发的微服务实现，主要实现都在`boot/service`包中，提供了一个基础的`BaseService`实现。因为使用`Etcd`服务注册和发现，所以需要依赖ETCD。

## 2. 基础概念

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


## 3. 实践

### 3.1 建立一个Service

#### 3.1.1 配置proto

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

#### 3.1.2 增加入口文件和Server

一个Service中，它的入口文件指向的是一个Server实现，这个Server是实现了pb中定义的服务接口。

server/server.go：

```go
package server

...

type (
	Server struct {
		pb.UnimplementedUserServiceServer
		service.BaseService
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

	// 此处关键
	pb.RegisterUserServiceServer(s.GrpcServer, s)

	return nil
}

```

同时，将入口文件中，是New此Server，来提供程序服务:

```go
package main

...

var (
	name       = flag.String("name", "UserService", "service name")
	addr       = flag.String("addr", "127.0.0.1:20001", "user service listen address")
	configPath = flag.String("config", "configs/user.service.yml", "config file path")
	instanceId = flag.Int64("instance", 1, "run instance id")
)

func main() {
	flag.Parse()

	instance := server.New(*name, *configPath, *addr, *instanceId)

	fmt.Println("server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}

```

#### 3.1.3 创建一个client

Service中，包含实现一个该服务的go语言client，提供给其他服务或API调用:

client/client.go

```
package client

...

// NewUserService 创建UserServcie客户端
func NewUserService(serviceName string, config *viper.Viper) (pb.UserServiceClient, error) {
	conn, err := service.NewGrpcClientConn(serviceName, config)
	if err != nil {
		return nil, err
	}
	return pb.NewUserServiceClient(conn), nil
}
```

#### 3.1.4 实现接口

server的方法中，必须实现pb定义接口，作为调用入口:

```go

// GetUserInfo 获取用户信息
func (s *Server) GetUserInfo(ctx context.Context, args *pb.GetUserInfoArgs) (*pb.GetUserInfoReply, error) {
	...
	return &pb.GetUserInfoReply{
		UserName: "hw",
	}, nil
}

```

### 3.1 逻辑处理

Service中，因为复杂性。因此要将逻辑，汇总到一个实现中，这个实现遵循DDD的限界上下文原则，需要的资源和数据，由Server在初始化其时，传递给它。此部分逻辑，建立在Service的internal目录之下：

```go

package logic

...

type (
	// InstanceLogic 实例的逻辑
	InstanceLogic struct {
		instanceManager *etcdx.ServiceManager

		instances map[string]Endpoint
		mu        sync.RWMutex
	}

	// Endpoint 实例节点
	Endpoint struct {
		Addr     string
		Metadata gameapi.GameMetaData
	}
)

// NewInstanceLogic 实例逻辑
func NewInstanceLogic(etcdUrl []string, instanceTarget string) (*InstanceLogic, error) {
	il := new(InstanceLogic)

	// 初始化InstanceManager
	xClient, err := etcdx.New(etcdUrl)
	if err != nil {
		return nil, err
	}

	il.instanceManager, err = etcdx.NewServiceManager(xClient, instanceTarget, 0, 0)
	if err != nil {
		return nil, err
	}

	// 加载节点数据
	m, _ := il.instanceManager.List()
	il.initInstances(m)

	// 监听数据变化
	il.instanceManager.Watch(il.updateInstace)

	return il, nil
}

```

> logic非必须，由业务的复杂度决定，如果只是简单CURD，可以直接调用Model


### 3.2 Model数据处理

对数据的处理，采取数据建模的设计。即，数据必然映射成某一个结构体。然后通过，Model定义一组对数据操作的接口，再根据使用的数据驱动（数据库），实现不同对数据的操作：

```go

package model

...

const (
	UserSkinDatabase  = "ET1"
	UserSkinTableName = "user_skin"
)

type (

	// UserSkin 玩家所选皮肤皮肤
	UserSkin struct {
		ID           int64            `bson:"_id" json:"id" dynamodbav:"id"`                        // 记录ID
		UserID       int64            `bson:"user_id" json:"user_id" dynamodbav:"user_id"`          // 用户ID
		SkinID       int64            `bson:"skin_id" json:"skin_id" dynamodbav:"skin_id"`          // 皮肤ID
		GSkinInfo    *pb.GSkinInfo    `bson:"skin_info" json:"skin_info" dynamodbav:"skin_info"`    // 皮肤详情，分自定义皮肤
		DiySkinParts map[string]int64 `bson:"skin_parts" json:"skin_parts" dynamodbav:"skin_parts"` // 自定义皮肤装扮
	}

	UserSkinModel interface {
		FindUpdate(ctx context.Context, data *UserSkin) error

		FindByUserID(ctx context.Context, userID int64) (*UserSkin, error)
	}

	defaultUserSkinModel struct {
		DB *mongo.Collection
	}
)

func NewUserSkinModel(mongoUrl string) UserSkinModel {
	collection := mongodbx.NewCollection(mongoUrl, UserSkinDatabase, UserSkinTableName)

	return &defaultUserSkinModel{
		DB: collection,
	}
}

func (m *defaultUserSkinModel) FindUpdate(ctx context.Context, data *UserSkin) error {

	...
}

// FindByUserID 用户ID查询指定用户选择皮肤
func (m *defaultUserSkinModel) FindByUserID(ctx context.Context, userID int64) (*UserSkin, error) {
	...
}
```

## 4. 注意事项