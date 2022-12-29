package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/trainking/goboot/internal/service/user/server"
)

var (
	name       = flag.String("name", "UserService", "service name")
	addr       = flag.String("addr", "127.0.0.1:20001", "user service listen address")
	configPath = flag.String("config", "configs/user.service.yml", "config file path")
	instanceId = flag.Int64("instance", 1, "run instance id")
)

func main() {
	flag.Parse()

	instance := server.New(*name, *configPath, *addr, *instanceId)

	// 1. 初始化
	if err := instance.Init(); err != nil {
		fmt.Println("server init failed, Error: ", err)
		return
	}

	// 2. 优雅退出
	go func() {
		exitC := make(chan os.Signal, 1)
		signal.Notify(exitC, syscall.SIGINT, syscall.SIGTERM)

		<-exitC
		instance.Stop()
		os.Exit(0)
	}()

	// 3. 启动实例
	fmt.Println("server start listen: ", *addr)
	if err := instance.Start(); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}
