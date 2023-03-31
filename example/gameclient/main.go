package main

import (
	"fmt"
	"time"

	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/gameapi"
	"github.com/xtaci/kcp-go"
	"google.golang.org/protobuf/proto"
)

func main() {
	c, e := kcp.Dial("127.0.0.1:6001")
	if nil != e {
		panic(e)
	}
	defer c.Close()

	go func() {
		for {
			n, e := gameapi.Packing(c)
			if nil != e {
				fmt.Println("read error:", e.Error())
				return
			}

			opecode := pb.OpCode(n.OpCode())
			switch opecode {
			case pb.OpCode_Pong:
				var result pb.S2C_Pong
				proto.Unmarshal(n.Body(), &result)

				fmt.Println(result.OK)
			}
		}
	}()

	msg := pb.C2S_Ping{TickTime: time.Now().Unix()}
	p, err := gameapi.CretaePbPacket(uint16(pb.OpCode_Ping), &msg)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := c.Write(p.Serialize()); err != nil {
		fmt.Println(err)
	}
	if _, err := c.Write(p.Serialize()); err != nil {
		fmt.Println(err)
	}

	time.Sleep(10 * time.Second)
	c.Close()
}
