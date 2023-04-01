package main

import (
	"fmt"
	"time"

	"github.com/trainking/goboot/example/gameclient/client"
	"github.com/trainking/goboot/internal/pb"
	"github.com/xtaci/kcp-go"
	"google.golang.org/protobuf/proto"
)

func main() {
	c, e := kcp.Dial("127.0.0.1:6001")
	if nil != e {
		panic(e)
	}
	defer c.Close()

	cc := client.NewClient(c, 1024, 1024, 3*time.Second)

	go func() {
		for p := range cc.Receive() {
			opcode := pb.OpCode(p.OpCode())
			switch opcode {
			case pb.OpCode_Op_S2C_Pong:
				var r pb.S2C_Pong
				proto.Unmarshal(p.Body(), &r)
				fmt.Printf("Pong: %v\n", r.OK)
			case pb.OpCode_Op_S2C_Login:
				var r pb.S2C_Login
				proto.Unmarshal(p.Body(), &r)
				fmt.Printf("Login: %v\n", r.Ok)
			}
		}
	}()

	msg := pb.C2S_Ping{TickTime: time.Now().Unix()}
	if err := cc.Send(uint16(pb.OpCode_Op_C2S_Ping), &msg); err != nil {
		fmt.Println(err)
		return
	}

	msgLogin := pb.C2S_Login{Account: "admin", Password: "123456"}
	if err := cc.Send(uint16(pb.OpCode_Op_C2S_Login), &msgLogin); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(2 * time.Second)
	if err := cc.Send(uint16(pb.OpCode_Op_C2S_Ping), &msg); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(6 * time.Second)
	if err := cc.Send(uint16(pb.OpCode_Op_C2S_Ping), &msg); err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(2 * time.Second)

	fmt.Println("end")
	cc.Close()
}
