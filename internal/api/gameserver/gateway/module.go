package gateway

import (
	"fmt"
	"strconv"

	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/gameapi"
	"github.com/trainking/goboot/pkg/log"
)

var Module = func() gameapi.Moddule {
	return new(GateWayM)
}

type GateWayM struct {
}

func (m *GateWayM) Init(a *gameapi.App) {

	a.SetConnectListener(func(s *gameapi.Session) error {
		log.Infof("ConnectNum: %d", a.GetTotalConn())
		return nil
	})

	a.SetDisconnectListener(func(s *gameapi.Session) error {
		log.Infof("ConnectNum: %d", a.GetTotalConn())
		return nil
	})

	// 设置消息处理前中间件
	a.SetBeforeMiddleware(func(s *gameapi.Session, p gameapi.Packet) error {
		if p.OpCode() == uint16(pb.OpCode_Op_C2S_Login) {
			return nil
		}

		if !s.IsValid() {
			return fmt.Errorf("session is valid, opcode: %d", p.OpCode())
		}
		return nil
	})
}

func (m *GateWayM) Group() map[uint16]gameapi.Handler {
	return map[uint16]gameapi.Handler{
		uint16(pb.OpCode_Op_C2S_Login): m.C2S_LoginHandler,
		uint16(pb.OpCode_Op_C2S_Say):   m.C2S_SayHandler,
	}
}

// C2S_LoginHandler 登录
func (m *GateWayM) C2S_LoginHandler(c gameapi.Context) error {
	var msg pb.C2S_Login
	if err := c.Params(&msg); err != nil {
		return err
	}

	log.Infof("Login: %s %s", msg.Account, msg.Password)
	var result = new(pb.S2C_Login)
	if msg.Password == "123456" {
		id, _ := strconv.ParseInt(msg.Account, 10, 64)
		c.Valid(id)
		result.Ok = true
	} else {
		result.Ok = false
	}

	return c.Send(uint16(pb.OpCode_Op_S2C_Login), result)
}

func (m *GateWayM) C2S_SayHandler(c gameapi.Context) error {
	var msg pb.C2S_Say
	if err := c.Params(&msg); err != nil {
		return err
	}

	log.Infof("Say %v", msg)

	return c.SendActor(msg.Actor, uint16(pb.OpCode_Op_S2C_Say), &pb.S2C_Say{
		Word: msg.Word,
	})
}
