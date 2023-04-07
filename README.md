# goboot

- [goboot](#goboot)
	- [概述](#概述)
	- [快速开始](#快速开始)
		- [开发一个http服务器](#开发一个http服务器)
		- [开发一个gRPC服务](#开发一个grpc服务)
		- [开发一个游戏服务器](#开发一个游戏服务器)
	- [惯例](#惯例)
	- [参考](#参考)


## 概述

**goboot**是个一个**Golang**开发过程中，使用的一些经验总结的库，涵盖`Web`开发和游戏开发领域。

## 快速开始

### 开发一个http服务器

**goboot**开发http服务器，只需要传入调用`httpapi`包，传入配置文件地址，监听端口，以及实例ID:

```go
instance := httpapi.New(*configPath, *addr, *instanceId)

...

// 启动
boot.BootServe(instance)
```

通过模块化的Module，来建立API，一个Module要遵循`httpapi.Module`的定义：

```golang
Module interface {
		// 初始化模块
		Init(app *App)
		// 模块的分组路由
		Group() Group
}
```

> 依赖nats做消息转发，Etcd做服务注册

### 开发一个gRPC服务

**goboot**开发gRPC服务器，同样只需要使用：

```
instance := server.New(*name, *configPath, *addr, *instanceId)

boot.BootServe(instance)
```

不同的是，这里的server需要定制protobuf文件实现。

### 开发一个游戏服务器

**goboot**开发游戏服务器，使用的`gameapi`包的实现：

```golang
instance := gameapi.New(*configPath, *addr, *instanceId)

boot.BootServe(instance)
```

同样使用的模块化的管理接口，不同的是，gameapi使用的`opcode->protobufMessage`的映射管理作为路由：

```go
Moddule interface {

		// 初始化模块
		Init(app *App)

		// 模块的分组路由
		Group() map[uint16]Handler
}
```

> 依赖Etcd做服务注册

## 惯例

**goboot**希望开发能够遵循一些管理来实现，这样可以大部分减少描述的篇幅，**遵循惯例，也方便多人开发的协作**。因此，推崇一下惯例：

1. 所有**实例**的配置文件都放在根目录下的`configs`目录
2. 所有**实例**的配置文件中，必须包含一份日志配置
```yaml
# 日志配置
Logger:
  # 日志输出级别，debug->wrong->error
  Level: "debug"
  # 日志文件的输出分类，文件名是 {target}_{instanceId}.log
  Target: "gameserver.api"
  # 日志输文件夹
  Outpath: "./logs"
```
3. 所有`api`接口的实现，都需要参考Module实现模块化
4. 所有`service`的实现，必须通过protobuff定义gRPC实现
5. 其他惯例，可以参考示例中的参考实现

## 参考

* [Game Server: 游戏服](./docs/guide/gameserver.md)
* [HTTP API: http开发Api接口](./docs/guide/httpapi.md)
* [Service: gRPC微服务](./docs/guide/service.md)
