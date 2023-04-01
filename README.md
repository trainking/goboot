# goboot

- [goboot](#goboot)
  - [概述](#概述)
  - [快速开始](#快速开始)
    - [开发一个http服务器](#开发一个http服务器)
    - [开发一个gRPC服务](#开发一个grpc服务)
    - [开发一个游戏服务器](#开发一个游戏服务器)
  - [参考](#参考)
    - [实例](#实例)
    - [Module](#module)
    - [Handler](#handler)
    - [服务注册和发现](#服务注册和发现)
    - [日志](#日志)
  - [注意事项](#注意事项)


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

## 参考

### 实例

### Module

### Handler

### 服务注册和发现

### 日志

## 注意事项

1. 使用`gameapi`开发游戏服务器时，opcode为0已经被心跳包占用