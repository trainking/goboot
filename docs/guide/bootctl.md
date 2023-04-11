# bootctl 命令行工具

- [bootctl 命令行工具](#bootctl-命令行工具)
  - [1. 概述](#1-概述)
  - [快速开始](#快速开始)
    - [安装](#安装)
    - [boot init：初始化项目](#boot-init初始化项目)


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
> --name 是必须的，项目名称；此名称最好与项目文件夹目录一致