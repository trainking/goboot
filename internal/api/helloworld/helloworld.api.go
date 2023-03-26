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

	// // 1. 初始化
	// if err := instance.Init(); err != nil {
	// 	log.Errorf("server init failed, Error: %v", err)
	// 	return
	// }

	// // 3. 优雅退出
	// go func() {
	// 	exitC := make(chan os.Signal, 1)
	// 	signal.Notify(exitC, syscall.SIGINT, syscall.SIGTERM)
	// 	<-exitC

	// 	instance.Stop()
	// 	os.Exit(0)
	// }()

	// // 4. 启动实例
	// log.Infof("server start listen: %s", *addr)
	// if err := instance.Start(); err != nil {
	// 	log.Errorf("server start failed, Error: %v", err)
	// 	return
	// }
}
