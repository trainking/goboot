package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/trainking/goboot/internal/api/helloworld/user"
	"github.com/trainking/goboot/pkg/httpapi"
	"github.com/trainking/goboot/pkg/log"
)

var (
	addr       = flag.String("addr", ":8001", "helloworld service listen address")
	configPath = flag.String("config", "configs/helloworld.api.yml", "config file path")
	instanceId = flag.Int64("instance", 1, "run instance id")
)

func main() {
	flag.Parse()

	instance := httpapi.New(*configPath, *addr, *instanceId)

	// 1. 初始化
	if err := instance.Init(); err != nil {
		log.Errorf("server init failed, Error: %v", err)
		return
	}

	// 2. 加载模块
	instance.AddModule(user.Module())

	// 3. 优雅退出
	go func() {
		exitC := make(chan os.Signal, 1)
		signal.Notify(exitC, syscall.SIGINT, syscall.SIGTERM)
		<-exitC

		instance.Stop()
		os.Exit(0)
	}()

	// 4. 启动实例
	log.Infof("server start listen: %s", *addr)
	if err := instance.Start(); err != nil {
		log.Errorf("server start failed, Error: %v", err)
		return
	}
}
