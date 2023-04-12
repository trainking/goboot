package main

import (
	"flag"
	"fmt"

	"github.com/trainking/goboot/internal/api/gameserver/gateway"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/gameapi"
)

var (
	name       = flag.String("name", "Gateway", "game server name")
	addr       = flag.String("addr", ":6001", "gameserver api lisen addr")
	configPath = flag.String("config", "configs/gameserver.api.yml", "config file path")
	instanceId = flag.Int64("instance", 1, "run instance id")
)

func main() {
	flag.Parse()

	instance := gameapi.New(*name, *configPath, *addr, *instanceId)

	instance.AddModule(gateway.Module())

	fmt.Println("game server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}
