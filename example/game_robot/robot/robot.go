package robot

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/gameapi"
	"github.com/xtaci/kcp-go"
	"google.golang.org/protobuf/proto"
)

type Robot struct {
	client *gameapi.Client
}

func New(network string, addr string) *Robot {
	r := new(Robot)

	var c net.Conn
	var e error
	switch network {
	case "kcp":
		c, e = kcp.Dial(addr)
	case "tcp":
		c, e = net.Dial("tcp", addr)
	}
	if nil != e {
		panic(e)
	}

	r.client = gameapi.NewClient(c, 1024, 1024, 3*time.Second)

	return r
}

// Login 玩家登录
func (r *Robot) Login(account string, password string) error {
	msg := pb.C2S_Login{Account: account, Password: password}
	if err := r.client.Send(uint16(pb.OpCode_Op_C2S_Login), &msg); err != nil {
		return err
	}

	p := <-r.client.Receive()
	if p.OpCode() == uint16(pb.OpCode_Op_S2C_Login) {
		var result pb.S2C_Login
		if err := proto.Unmarshal(p.Body(), &result); err != nil {
			return err
		}

		if !result.Ok {
			return errors.New("account or password is wrong")
		}
	}

	return nil
}

func (r *Robot) Receive() {
	defer func() {
		r.Quit()
	}()
	for p := range r.client.Receive() {
		opcode := pb.OpCode(p.OpCode())
		switch opcode {
		case pb.OpCode_Op_S2C_Say:
			var _resultMSg pb.S2C_Say
			if err := proto.Unmarshal(p.Body(), &_resultMSg); err != nil {
				fmt.Printf("%v\n", err)
				return
			}

			fmt.Printf("Say %v\n", _resultMSg)
		}
	}
}

// Say 向其他玩家发送消息
func (r *Robot) Say(actor int64, word string) error {
	msg := pb.C2S_Say{Actor: actor, Word: word}
	if err := r.client.Send(uint16(pb.OpCode_Op_C2S_Say), &msg); err != nil {
		return err
	}

	return nil
}

// Quit 机器人退出
func (r *Robot) Quit() {
	r.client.Close()
}
