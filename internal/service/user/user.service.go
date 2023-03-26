package main

import (
	"flag"
	"fmt"

	"github.com/trainking/goboot/internal/service/user/server"
	"github.com/trainking/goboot/pkg/boot"
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

	fmt.Println("server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}
