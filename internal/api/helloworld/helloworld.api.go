package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v4/middleware"
	"github.com/trainking/goboot/internal/api/helloworld/user"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/httpapi"
)

var (
	addr       = flag.String("addr", ":8001", "helloworld service listen address")
	configPath = flag.String("config", "configs/helloworld.api.yml", "config file path")
	instanceId = flag.Int64("instance", 1, "run instance id")
)

func main() {
	flag.Parse()

	instance := httpapi.New(*configPath, *addr, *instanceId)
	// 中间件
	instance.Use(middleware.RequestID())
	// 模块
	instance.AddModule(user.Module())

	fmt.Println("server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}
