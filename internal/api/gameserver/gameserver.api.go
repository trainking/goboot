package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/gameapi"
	"github.com/trainking/goboot/pkg/log"
)

var (
	addr       = flag.String("addr", ":6001", "gameserver api lisen addr")
	configPath = flag.String("config", "configs/gameserver.api.yml", "config file path")
	instanceId = flag.Int64("instance", 1, "run instance id")
)

func main() {
	flag.Parse()

	instance := gameapi.New(*configPath, *addr, *instanceId)

	instance.AddHandler(uint16(pb.OpCode_Ping), HandlerPing)

	fmt.Println("game server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}

// HandlerPing
func HandlerPing(session *gameapi.Session, b []byte) error {

	var msg pb.C2S_Ping
	proto.Unmarshal(b, &msg)

	log.Infof("Receive %v", msg.TickTime)
	time.Sleep(3 * time.Second)
	return nil
}
