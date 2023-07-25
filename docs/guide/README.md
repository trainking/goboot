# GOBOOT 指南

- [GOBOOT 指南](#goboot-指南)
  - [概述](#概述)
  - [组成](#组成)
  - [连接](#连接)
  - [服务](#服务)
  - [最佳实践](#最佳实践)
  - [惯例](#惯例)
  - [参考文件](#参考文件)


## 概述

`GOBOOT`是一个golang开发的游戏服开发框架。可以用于开发，游戏周边系统，游戏后台系统（无物理部分）。也可以用于开发其他应用系统。实现了`HTTP`，`TCP`，`KCP`, `Websocket`协议的接入。

`GOBOOT`使用分层架构设计，同时借鉴`DDD`(领域驱动设计)的优秀设计。遵循`连接-逻辑-数据`分层，同时根据业务划分成各自领域，每个领域之间通过通信机制交互。

## 组成

## 连接

`GOBOOT`通过建立API，来开放访问的连接。存在两种API，分别是`httpapi`和`gameapi`。

* httpapi: 接受http连接，短链接服务实现。
* gameapi: 接受`TCP`，`KCP`，`Websocket`，长链接服务实现。

## 服务

`GOBOOT`将逻辑和数据合并成`Service`，因为逻辑本省依赖着数据。同时，每一个`Servcie`可以看作一个领域，因此通过构造实现的隔离。将逻辑需要的数据，通过注入到上下文的方式，提供给逻辑使用。

## 最佳实践

## 惯例

开发能够遵循一些管理来实现，这样可以大部分减少描述的篇幅，**遵循惯例，也方便多人开发的协作**。因此，推崇一下惯例：

1. 所有**实例**的配置文件都放在根目录下的`configs`目录
2. 所有**实例**的配置文件中，设定日志配置条件
```yaml
# 日志配置
Logger:
  # 日志输出级别，debug->info->warn->error
  Level: "debug"
  # 日志文件的输出分类，文件名是 {target}_{instanceId}.log
  Target: "gameserver.api"
  # 日志输文件夹
  Outpath: "./logs"
```
1. 所有`api`接口的实现，都需要参考Module实现模块化
2. 所有`service`的实现，必须通过protobuff定义gRPC实现
3. 其他惯例，可以参考示例中的参考实现

## 参考文件

* [echo文档](https://echo.labstack.com/docs/category/guide)
