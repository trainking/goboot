# bootctl 命令行工具

- [bootctl 命令行工具](#bootctl-命令行工具)
  - [1. 概述](#1-概述)
  - [快速开始](#快速开始)
    - [安装](#安装)
    - [boot init：初始化项目](#boot-init初始化项目)
    - [boot http: 创建一个http api](#boot-http-创建一个http-api)
    - [boot game：创建一个game api](#boot-game创建一个game-api)
    - [boot service：创建一个service](#boot-service创建一个service)


## 1. 概述

**bootctl**是goboot的命令行工具，它提供项目代码的生成和管理，方便快捷的生成项目结构，遵循惯例，减少新人的试错。

## 快速开始

### 安装

安装只需要使用`go install`即可:

```bash
> cd ./cmd/bootctl/
> go install .
```

### boot init：初始化项目

```
> bootctl init --name=project1
```

| 参数 | 类型   | 参考   | 是否必须 | 说明                                                         |
| :--- | :----- | :----- | :------- | :----------------------------------------------------------- |
| name | string | auther | 是       | 项目名称，也是go module 名称；此名称最好与项目文件夹目录一致 |

### boot http: 创建一个http api

```
> bootctl http --name=hello
```

| 参数 | 类型   | 参考           | 是否必须 | 说明                 |
| :--- | :----- | :------------- | :------- | :------------------- |
| name | string | auther         | 是       | 服务名               |
| addr | stirng | 127.0.0.1:8080 | 否       | 默认监听的地址和端口 |
| id   | int    | 1              | 否       | 默认的实例id         |

### boot game：创建一个game api

```
> bootctl init --name=gateway
```

| 参数    | 类型   | 参考                  | 是否必须 | 说明                                                                     |
| :------ | :----- | :-------------------- | :------- | :----------------------------------------------------------------------- |
| name    | string | auther                | 是       | 服务名                                                                   |
| option  | string | c                     | 否       | 操作符，默认是c; c 创建；g 生成proto，会根据编辑的game.proto生成op.proto |
| addr    | stirng | 127.0.0.1:8080        | 否       | 默认监听的地址和端口                                                     |
| id      | int    | 1                     | 否       | 默认的实例id                                                             |
| network | string | kcp                   | 否       | 传输协议，tcp或kcp，默认kcp                                              |
| nats    | string | nats://127.0.0.1:4222 | 否       | nats地址                                                                 |
| etcd    | stirng | 127.0.0.1:2379        | 否       | etcd地址，集群用`,`分割                                                  |

编辑xxx.game.proto后，生成xxx.op.proto:

```
bootctl init --name=gateway --option=g
```

### boot service：创建一个service

```
> bootctl service --name=user
```

| 参数 | 类型   | 参考           | 是否必须 | 说明                 |
| :--- | :----- | :------------- | :------- | :------------------- |
| name | string | auther         | 是       | 服务名               |
| addr | stirng | 127.0.0.1:8080 | 否       | 默认监听的地址和端口 |
| id   | int    | 1              | 否       | 默认的实例id         |