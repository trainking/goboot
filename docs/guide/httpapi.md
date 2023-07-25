# HTTP API 的最佳实践

- [HTTP API 的最佳实践](#http-api-的最佳实践)
	- [1. 概述](#1-概述)
	- [2. 基础概念](#2-基础概念)
		- [2.1 Instance](#21-instance)
			- [2.1.1 配置文件](#211-配置文件)
		- [2.2 Module](#22-module)
			- [2.2.1 Group](#221-group)
			- [2.2.2 Router](#222-router)
		- [2.3 Handler](#23-handler)
	- [3. 实践](#3-实践)
		- [3.1 建立一个http实例](#31-建立一个http实例)
		- [3.2 API如何访问Service](#32-api如何访问service)
	- [4. 注意事项](#4-注意事项)


## 1. 概述

**goboot**开发基于HTTP协议的API接口，使用`pkg/httpapi`包，目前只实现了`HTTP 1.1`。这个包是通过开源框架`echo`定制开发，所以很多内容可以查看:
* [echo: High performance, minimalist Go web framework](https://github.com/labstack/echo)

## 2. 基础概念

### 2.1 Instance

#### 2.1.1 配置文件

`http api`的配置都是根据实际业务可选配置。

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

## 3. 实践

### 3.1 建立一个http实例

1. 在根目录下的`internal/api`目录下，增加一个具体实例名的文件夹，例如`user`
2. 此包中建立一个`xxx.api.go`的入口文件，如`user.api.go`, 范例如下:

```go
package main

...

var (
	name       = flag.String("name", "helloworld", "http api server name")
	addr       = flag.String("addr", ":8001", "helloworld service listen address")
	configPath = flag.String("config", "configs/helloworld.api.yml", "config file path")
	instanceId = flag.Int64("instance", 1, "run instance id")
)

func main() {
	flag.Parse()

	instance := httpapi.New(*name, *configPath, *addr, *instanceId)
	
	// 中间件

	// 模块
	instance.AddModule(user.Module())

	fmt.Println("server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}
```

3. 根据具体业务建立模块，在此包中，建立具体模块包，同时模块包中，约定必有一个`module.go`文件，此文件中初始化一个模块：

```go
package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/httpapi"
	"github.com/trainking/goboot/pkg/log"
	"github.com/trainking/goboot/pkg/service"

	userClient "github.com/trainking/goboot/internal/service/user/client"
)

var Module = func() httpapi.Module {
	return new(M)
}

// M 定义模块具有的上下文
type M struct {
	Config *viper.Viper

}

// Init 初始化模块必须的资源
func (m *M) Init(app *httpapi.App) {

	log.Debugf("Init Module Config: %v", app.Config)
	m.Config = app.Config

	...
}

// Group 设定访问路径
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
}

```

4. 建立一个handler，监听配置在`Group`中:

```
func (m *M) GetUserInfo(c echo.Context) error {
	requestID := httpapi.GetRequestID(c)
	log.Trace(requestID, "GetUserInfo", "start", "start GetUserInfo")
	ctx := service.WithRequestIDContext(c.Request().Context(), requestID)
	reply, err := m.UserService.GetUserInfo(ctx, &pb.GetUserInfoArgs{
		UserId: 1,
	})
	if err != nil {
		log.Errorf("GetUserInfo error: %v", err)
	}
	return c.JSON(http.StatusOK, reply.UserName)
}
```

5. 将Module在`xxx.api.go`文件中，通过`AddModule`函数加入到实例中，便可以应用。注意，要注意加入的位置，应在New和Start之间。

### 3.2 API如何访问Service

API本身不实现具体业务，只关心连接有关的。具体的业务实现，都是交给`Service`实现。因此，API调用Service，需要满足一下条件：

1. 在API配置文件，加入对应Service注册到Etcd的信息：

```yaml
UserService:
  Prefix: "/services"
  Authentication:
    Mode: "TLS"
    CertFile: "../ssldata/cert.crt"
  Etcd:
    - "192.168.1.9:2379"
```

2. 在Module的Init方法，将此客户端初始化：

```go
	var err error
	m.UserService, err = userClient.NewUserService("UserService", app.Config)
	if err != nil {
		log.Errorf("userClient.NewUserService failed: %v", err)
		return
	}
```

## 4. 注意事项

> httpapi没有做服务注册，因为http由很多负载均衡组件，不需要过多的定制化