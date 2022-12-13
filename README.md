# goboot

- [goboot](#goboot)
  - [概述](#概述)
  - [开发指南](#开发指南)
    - [目录](#目录)
    - [设计](#设计)


## 概述

**goboot**是个一个**Golang**开发过程中，使用的一些经验总结的库，涵盖`Web`开发和游戏开发领域。

## 开发指南

### 目录

* `bin`: 编译后程序目录
* `cmd`: 自定义命令工具
* `configs`: 配置文件存放目录
* `internal`: 具体项目的逻辑存放目录
* `pkg`: 通用包目录
  - `boot`：项目启动设定包
  - `encrypt`：加密解密包
  - `errgroup`：一个携程编组实现
  - `etcdx`: etcd操作拓展包
  - `httpapi`: 基于Echo框架实现的http api实现
  - `idgen`：id生成器包
  - `jwt`： jwt的实现
  - `log`: 日志库的实现
  - `mongodbx`: mongodb操作拓展包
  - `rabbitmqx`：rabbitmq操作拓展包
  - `random`: 随机函数包
  - `redisx`: 基于redis6.0的工具实现
  - `utils`: 一些工具函数的实现

### 设计

该脚手架使用多框架融合，主要使用以下：

* Echo: 一个http的api框架
* gRPC：goole开源的RPC框架

具体使用，参考如下文档：

1. [设计思想](./docs/guide/%E8%AE%BE%E8%AE%A1%E6%80%9D%E6%83%B3.md)