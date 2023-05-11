# Game Server: 游戏服

- [Game Server: 游戏服](#game-server-游戏服)
  - [1. 概述](#1-概述)
  - [1.1 Actor模式](#11-actor模式)
  - [2. 详细设计](#2-详细设计)
    - [2.1 协议](#21-协议)
    - [2.2 Instance](#22-instance)
      - [2.2.1 配置文件](#221-配置文件)
      - [2.2.2 应用协议](#222-应用协议)
    - [2.3 Module](#23-module)
    - [2.4 Handler](#24-handler)
      - [SendXXX](#sendxxx)
  - [3. 注意事项](#3-注意事项)


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

## 2. 详细设计

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

## 3. 注意事项

> 使用`gameapi`开发游戏服务器时，opcode为0已经被心跳包占用