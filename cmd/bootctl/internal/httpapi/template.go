package httpapi

var mainText = `package main

import (
	"flag"
	"fmt"

	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/httpapi"
)

var (
	name       = flag.String("name", "{{name}}", "http api server name")
	addr       = flag.String("addr", "{{addr}}", "{{name}} http api listen address")
	configPath = flag.String("config", "configs/{{name}}.api.yml", "config file path")
	instanceId = flag.Int64("instance", {{id}}, "run instance id")
)

func main() {
	flag.Parse()

	instance := httpapi.New(*name, *configPath, *addr, *instanceId)
	
	// TODO  增加中间价
	
	// TODO 增加Module

	fmt.Println("server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}`
