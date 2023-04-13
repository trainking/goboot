# HTTP API

- [HTTP API](#http-api)
	- [1. 概述](#1-概述)
	- [2. 详细设计](#2-详细设计)
		- [2.1 Instance](#21-instance)
			- [2.1.1 配置文件](#211-配置文件)
		- [2.2 Module](#22-module)
			- [2.2.1 Group](#221-group)
			- [2.2.2 Router](#222-router)
		- [2.3 Handler](#23-handler)
	- [3. 注意事项](#3-注意事项)


## 1. 概述

**goboot**开发基于HTTP协议的API接口，使用`pkg/httpapi`包，目前只实现了`HTTP 1.1`。这个包是通过开源框架`echo`定制开发，所以很多内容可以查看:
* [echo: High performance, minimalist Go web framework](https://github.com/labstack/echo)

## 2. 详细设计

### 2.1 Instance

`Instance`即实例，是一个http服务器的最小执行单位，可以看作是一个进程实例。它必须传入三个参数:

| 参数       | Options               | 类型   | 说明                              |
| :--------- | :-------------------- | :----- | :-------------------------------- |
| name       | -name,--name          | string | 名称                              |
| addr       | -addr, --addr         | string | 监听的ip:port，如`127.0.0.1:6001` |
| configPath | -config, --config     | string | 配置文件的路径                    |
| instanceId | -instance, --instance | int64  | 实例的ID，实例的唯一标识          |

> httpapi没有做服务注册，因为http由很多负载均衡组件，不需要过多的定制化

#### 2.1.1 配置文件

除了统一的日志配置文件之外，无需其他配置

### 2.2 Module

`Module`是用来组合各个接口的Handler，其实现必须实现`httpapi`中的`Module`接口:

```go
// Module 按模块组合
	Module interface {
		// 初始化模块
		Init(app *App)
		// 模块的分组路由
		Group() Group
	}
```

#### 2.2.1 Group

`Group`是自定义的路由组，这里建立的一个关系，即**一个Module对应一个Group**。所有，路由都是定义在`Module`的`Group()`方法中:

```go
func (m *M) Group() httpapi.Group {
	return httpapi.Group{
		Path:        "/api/helloworld",
		Middlewares: []httpapi.Middleware{},
		Routers: []httpapi.Router{
			{
				Method: http.MethodGet,
				Path:   "/userinfo",
				Handle: m.GetUserInfo,
			},
		},
}
```


#### 2.2.2 Router

`Router`是自定义的路由封装，是对echo的包的二次抽象：

```
// Router 路由
Router struct {
    Method      string           // 方法
    Path        string           // 路径
    Name        string           // 名称
    Handle      echo.HandlerFunc // 处理函数
    Middlewares []Middleware     // 中间件函数
}
```

可以通过Instance的`AddRouter`直接加入顶级路由：

```
instance.AddRouter(httpapi.Router{...})
```

### 2.3 Handler

`Handler`就是使用的原有的`echo#Handler`，未做定制化修改。

## 3. 注意事项