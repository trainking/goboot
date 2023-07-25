# Game API的最佳实践

- [Game API的最佳实践](#game-api的最佳实践)
  - [1. 概述](#1-概述)
  - [1.1 Actor模式](#11-actor模式)
  - [2. 基础概念](#2-基础概念)
    - [2.1 协议](#21-协议)
    - [2.2 Instance](#22-instance)
      - [2.2.1 配置文件](#221-配置文件)
      - [2.2.2 应用协议](#222-应用协议)
    - [2.3 Module](#23-module)
    - [2.4 Handler](#24-handler)
      - [SendXXX](#sendxxx)
    - [2.5 Listener](#25-listener)
    - [2.6 Middleware](#26-middleware)
  - [3. 实践](#3-实践)
    - [3.1 建立一个game api的实例](#31-建立一个game-api的实例)
    - [3.2 增加建立连接监听和断开连接监听](#32-增加建立连接监听和断开连接监听)
    - [3.3 增加handler处理的前置和后置条件](#33-增加handler处理的前置和后置条件)
  - [4. 注意事项](#4-注意事项)


## 1. 概述

**goboot**开发游戏服，使用的是`pkg/gameapi`包。这个包中，定义了传输协议，序列化结构，以及服务的注册发现。其中，启动一个游戏服，必须要使用的两个依赖，这两个组件，都需要在实例的配置文件中配置：

* nats: 高性能的消费分发组件，用于在多实例中发布/订阅消息
* Etcd：服务注册发现的组件

> 这里的游戏服，特指为游戏服务的服务器，即长链接服务器。如需要，开发http请求，可以通过httpapi包，定义一套配套API

## 1.1 Actor模式

> Actor模式是一种并发计算模型，旨在提高并发程序的可伸缩性和容错性。它最初由Carl Hewitt于1973年提出，并被称为Actor模型。
> 在Actor模型中，计算单元被称为Actor。每个Actor都是独立的，具有自己的状态和行为，并且与其他Actor通过异步消息传递进行通信。

这是标准的Actor模型，但是标准的Actor模型过于灵活，且对于程序的理解存在的很大问题。因此，goboot借鉴了Actor一些思想，采用的是`session-handler`方式。

其中，`Session`可以看作是一个Actor，它的信箱由两个`chan`组成：

```golang
type Session struct {
  ...

  sendChan    chan Packet // 发送队列
	receiveChan chan Packet // 接收队列

  ...
}
```

同时，对于message的处理，则是通过判断是否由注册Handler，如果由注册Handler，则处理，没有则直接**发送会客户端**；这样设计就意味着，对于客户端发送的消息，**必须要Handler**。

使用Nats作为消息总线（也可以说是`Actor Push`），可以很好保证并发。

## 2. 基础概念

### 2.1 协议

`gameapi` 自定义包协议来实现数据传输，传输层协议可以使用`tcp`和`kcp`。包协议使用了**定长包头+变长包体**方案：

```
 |-----------------------------message-----------------------------------------|
 |----------------------Header------------------|------------Body--------------|
 |------Body Length-------|--------Opcode-------|------------Body--------------|
 |---------uint16---------|---------uint16------|------------bytes-------------|
 |-----------2------------|----------2----------|-----------len(Body)----------|
```

* 头两个字节保存body的字节大小
* 后两个字节保存body对应的Opcode，用于索引对应的处理Handler可以解包结构
* 网络字节序，使用**大端序**

### 2.2 Instance

`Instance`是游戏服最小的执行单位，可以粗略认为的，这就是一个游戏进程。一个`Instance`需要传入四个必须参数:

| 参数       | Options               | 类型   | 说明                                                                          |
| :--------- | :-------------------- | :----- | :---------------------------------------------------------------------------- |
| name       | -name,--name          | string | 名称，定义此实例作用，可以根据name组合Module实现不同作用的实例，如Gate, Lobby |
| addr       | -addr, --addr         | string | 监听的ip:port，如`127.0.0.1:6001`                                             |
| configPath | -config, --config     | string | 配置文件的路径                                                                |
| instanceId | -instance, --instance | int64  | 实例的ID，实例的唯一标识                                                      |

#### 2.2.1 配置文件

示例：

```yaml
# 传输层协议，tcp, kcp
Network: "kcp"
# 每个连接的读超时(等于客户端心跳的超时)，秒为单位
ConnReadTimeout: 10
# 每个连接的写超时，秒为单位
ConnWriteTimeout: 5
# 连接成功后，多久未验证身份，则断开，秒为单位
ValidTimeout: 10
# 最大发送消息包大小
SendLimit: 1024
# 最大接收消息包大小
ReceiveLimit: 1024
# 心跳包限制数量, 每分钟不能超过的数量
HeartLimit: 100
# NATS的地址
NatsUrl: "nats://192.168.1.14:4222"
# 服务注册的前缀
Prefix: "/gameserver"
# 服务注册的Etcd地址
Etcd:
  - "192.168.1.9:2379"
```

配置详解：

| 配置项          | 类型              | 是否必须 | 说明                                                    |
| :-------------- | :---------------- | :------- | :------------------------------------------------------ |
| Prefix          | string            | 是       | 服务注册的前缀，服务注册key格式: <prefix>/<name>/<addr> |
| Etcd            | []string          | 是       | Etcd地址，服务注册需要                                  |
| NatsUrl         | string            | 是       | NATS的地址                                              |
| Network         | string            | 是       | 传输协议，可以是`tcp`, `kcp`, `websocket`               |
| WebsocketPath   | string            | 否       | 传输协议是`websocket`时必须                             |
| KcpMode         | string            | 否       | kcp模式，nomarl 普通模式 fast 极速模式；默认极速模式    |
| ConnReadTimeout | int               | 是       | 每个连接的读超时(等于客户端心跳的超时)，秒为单位        |
| ConnReadTimeout | int               | 是       | 每个连接的写超时，秒为单位                              |
| ValidTimeout    | int               | 是       | 连接成功后，验证身份超时，秒为单位                      |
| SendLimit       | int               | 是       | 最大发送消息包大小；发送消息缓冲区大小                  |
| ReceiveLimit    | int               | 是       | 最大接收消息包大小；接收消息缓冲区大小                  |
| HeartLimit      | int               | 是       | 心跳包限制数量, 每分钟不能超过的个数                    |
| Password        | string            | 否       | 服务加密使用的密码                                      |
| TLS             | map[string]string | 否       | TLS配置，配置了才开启                                   |
| TLS.CertFile    | string            | 否       | cert文件路径                                            |
| TLS.KeyFile     | string            | 否       | 密钥文件路径                                            |

#### 2.2.2 应用协议

`gameapi`的应用协议，推荐使用`protobuf`作为传输协议定义，推荐如下结构:

```protobuf
syntax = "proto3";

option go_package = "./pb";

package pb;

// OpCode 操作符定义
enum OpCode {
  None = 0;  // 心跳包占用
  Op_C2S_Login = 1;
  Op_S2C_Login = 2;
  Op_C2S_Say = 3;
  Op_S2C_Say = 4;
}

// C2S_Login 登录服务器
message C2S_Login {
  string Account = 1;
  string Password = 2;
}
message S2C_Login { bool Ok = 1; }

// C2S_Say 玩家通信
message C2S_Say {
  int64 Actor = 1;
  string Word = 2;
}
message S2C_Say {
  string Word = 1;
}
```

### 2.3 Module

`Module`是实例中相关内容的集合，可以认为实例就是由一个个Module组成。`gameapi`中的模块必须实现`gameapi.Module`：

```go
// Module 模块
	Moddule interface {

		// 初始化模块
		Init(app *App)

		// 模块的分组路由
		Group() map[uint16]Handler
	}
```

开发过程中，遵循，在Init函数中初始化Handler需要的各项资源，在Group中定义该模块的`opcode->handler`的映射关系。

### 2.4 Handler

`Handler`是对消息处理的单元，一个Handler对应一个消息的处理，实现一个`Handler`如下:

```
// Handler最好是Module的函数，方便引用资源
func (m *GateWayM) C2S_SayHandler(c gameapi.Context) error {
	
    // 读取输入消息
    var msg pb.C2S_Say
	if err := c.Params(&msg); err != nil {
		return err
	}

    // 向该玩家返回消息
	return c.SendActor(msg.Actor, uint16(pb.OpCode_Op_S2C_Say), &pb.S2C_Say{
		Word: msg.Word,
	})
}
```

#### SendXXX

在Handler中，可以使用多种`ctx.SendXXX`方法，向玩家和`Actor push`发送消息:

* Send：发送消息到玩家客户端
* SendActor: 向指定玩家发送Actor消息，如果有Handler则会被Handler处理；如无，则会被直接发给玩家。同时，这个方法中，会优先判断玩家是否在本实例，在本实例，则直接发送，不会发送`Actor push`。
* SendAllActor：向所有玩家广播一个Actor
* SendActorLocation：发送到本实例指定玩家的Actor
* SendActorPush: 发送消息到`Actor push`

### 2.5 Listener

Listener是监听连接状态的监听器。`gamapi`中提供两个监听器，分别是`ConnectListener`和`DisconnectListener`，分别代表对连接建立和连接断开的监听。可以通过`app.SetConnectListener()`和`app.SetDisconnectListener`设置它们的实现，它们的实现必须是一个Listener函数:

```go
// Listener 监听器
	Listener func(*Session) error
```

> 注意：这两个监听器在gameapi中，保持单个实现，重复设置将会被覆盖

### 2.6 Middleware

Middleware是处理消息时的中间件，`gameapi`提供了两种中间件，分别是`beforeMiddleware`和`afterMiddleware`两种，分别代表前置中间件和后置中间件，发生在处理消息之前和处理之后。它们都可以定义多个，每个都会通过`Condition`条件，判断是否可以被执行:

```go
  // Middleware 中间件处理
	Middleware struct {
		// Condition 是否要处理的opcode
		Condition func(uint16) bool
		// Do 处理执行
		Do func(Context) error
	}
```

## 3. 实践

### 3.1 建立一个game api的实例

1. 在根目录下的`internal/api`目录下，增加一个具体实例名的文件
2. 此包中建立一个`xxx.api.go`的入口文件， 范例如下:

```go
package main

...

var (
	name       = flag.String("name", "Gateway", "game server name")
	addr       = flag.String("addr", ":6001", "gameserver api lisen addr")
	configPath = flag.String("config", "configs/gameserver.api.yml", "config file path")
	instanceId = flag.Int64("instance", 1, "run instance id")
)

func main() {
	flag.Parse()

	instance := gameapi.New(*name, *configPath, *addr, *instanceId)

  // 增加模块
	instance.AddModule(gateway.Module())

	fmt.Println("game server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}
```

3. 根据具体业务建立模块，在此包中，建立具体模块包，同时模块包中，约定必有一个`module.go`文件，此文件中初始化一个模块：

```go
var Module = func() gameapi.Module {
	return new(GateWayM)
}

type GateWayM struct {
}

func (m *GateWayM) Init(a *gameapi.App) {
	log.Info("Module init")

	...

  // 增加opcode与handler处理映射
	a.AddHandler(pb.OpCode_Op_C2S_Login, m.C2S_LoginHandler)
	a.AddHandler(pb.OpCode_Op_C2S_Say, m.C2S_SayHandler)
	a.AddHandler(pb.OpCode_Op_S2S_Hi, m.S2S_Hi)
}
```

4. 将Module在`xxx.api.go`文件中，通过`AddModule`函数加入到实例中，便可以应用。注意，要注意加入的位置，应在New和Start之间。

### 3.2 增加建立连接监听和断开连接监听

建立连接和断开连接，是一个长连接服务器中，非常重要的事件。监听这两个事件，我们可以在连接建立时，做一些初始化的操作，在连接断开时，做一些销毁动作。`GOBOOT`中，增加了两个函数，分别是:

```go

  // 建立连接监听事件
  a.SetConnectListener(func(s *gameapi.Session) error {
		log.Infof("ConnectNum: %d", a.GetTotalConn())
		return nil
	})

  // 断开连接的监听事件
	a.SetDisconnectListener(func(s *gameapi.Session) error {
		log.Infof("ConnectNum: %d", a.GetTotalConn())
		return nil
	})
```

### 3.3 增加handler处理的前置和后置条件

消息处理的前置和后置，会在消息处理之前和消息处理之后，执行处理。需要注意的是，必须满足`Condition`条件，Do函数才会被执行。

```go
  // 设置消息处理前中间件
	a.AddBeforeMiddleware(gameapi.Middleware{
		Condition: func(opcode uint16) bool {
			return opcode != uint16(pb.OpCode_Op_C2S_Login)
		},
		Do: func(ctx gameapi.Context) error {
			if !ctx.Session().IsValid() {
				return fmt.Errorf("session is valid, opcode: %d", ctx.GetOpCode())
			}
			return nil
		},
	})
```

## 4. 注意事项

> 使用`gameapi`开发游戏服务器时，opcode为0已经被心跳包占用