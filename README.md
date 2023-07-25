# GOBOOT

## 概述

**GOBOOT** 是一个可以用来开http，gRPC和自定义游戏服务器的`Golang`脚手架。

## 特性

* 高性能分布式服务器开发
* 支持gRPC开发微服务
* 支持http的API开发
* 支持游戏服务器网关开发，支持`tcp`，`kcp`，`websocket`协议
* 提供了`bootctl`作为，项目管理和代码生成工具

## 快速开始

### 使用bootctl

**goboot**推荐使用`bootctl`命令行工具，来生成代码结构和代码，方便在多人开发时，控制编码规范。具体[参考](./docs/guide/bootctl.md)。

一个规范的目录结构如下:

```
├─bin                             # build之后的发布包
├─cmd                             # 自定义的命令行工具
├─configs                         # 所有的配置文件
├─docs                            # 项目相关文档
|  ├─devel                        # 开发相关文档，如接口描述
│  └─guide                        # 教程相关文档
├─example                         # 一些小示例，或测试demo
├─internal                        # 业务代码
│  ├─api                          # 应用程序接口，http和game
│  │  ├─gameserver
│  │  │  └─gateway
│  │  └─helloworld
│  │      └─user
│  ├─pb                           # protbuf生成库
│  │  └─proto                     # .proto文件
│  └─service                      # 服务
│      └─user
│          ├─client
│          └─server
├─logs                            # 日志
│  └─Gateway
└─pkg                             # 公共包
```

> 详情请查阅[参考指南](./docs/guide/README.md)


## 关联项目

* [goboot-csharp-client](https://github.com/trainking/goboot-csharp-client): gamapi的C#客户端实现，可用于unity开发

## 打赏作者

![打赏](./docs/guide/image/w_20230424115445.png)
