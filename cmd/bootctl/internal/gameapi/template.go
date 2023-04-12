package gameapi

var mainText = `package main

import (
	"flag"
	"fmt"

	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/gameapi"
)

var (
	name       = flag.String("name", "{{name}}", "game server name")
	addr       = flag.String("addr", "{{addr}}", "{{name}} game server lisen addr")
	configPath = flag.String("config", "configs/{{name}}.api.yml", "config file path")
	instanceId = flag.Int64("instance", {{id}}, "run instance id")
)

func main() {
	flag.Parse()

	instance := gameapi.New(*name, *configPath, *addr, *instanceId)

	// TODO add Module

	fmt.Println("game server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}`

var ymlText = `# 传输层协议，tcp, kcp
Network: "{{network}}"
# 每个连接的读超时(等于客户端心跳的超时)，秒为单位
ConnReadTimeout: 10
# 每个连接的写超时，秒为单位
ConnWriteTimeout: 5
# 连接成功后，多久未验证身份，则断开，秒为单位
ValidTimeout: 10
# 最大发送消息包大小
SendLimit: 1024
# 最大接收消息包大小
ReceiveLimit: 1024
# 心跳包限制数量, 每分钟不能超过的数量
HeartLimit: 100
# NATS的地址
NatsUrl: "{{nats}}"
# 服务注册的前缀
Prefix: "/gameserver"
# 服务注册的Etcd地址
Etcd:
{{etcd}}`