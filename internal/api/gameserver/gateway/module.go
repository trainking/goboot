package gateway

import (
	"errors"
	"fmt"

	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/gameapi"
	"github.com/trainking/goboot/pkg/log"
	"google.golang.org/protobuf/proto"
)

var Module = func() gameapi.Moddule {
	return new(GateWayM)
}

type GateWayM struct {
}

func (m *GateWayM) Init(a *gameapi.App) {

	a.SetBeforeMiddleware(func(s *gameapi.Session, p gameapi.Packet) error {
		if p.OpCode() == uint16(pb.OpCode_Op_C2S_Ping) || p.OpCode() == uint16(pb.OpCode_Op_C2S_Login) {
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
		uint16(pb.OpCode_Op_C2S_Ping):  m.C2S_PingHandler,
		uint16(pb.OpCode_Op_C2S_Login): m.C2S_LoginHandler,
	}
}

// C2S_PingHandler
func (m *GateWayM) C2S_PingHandler(session *gameapi.Session, b []byte) error {

	var msg pb.C2S_Ping
	proto.Unmarshal(b, &msg)

	log.Infof("Receive %v", msg.TickTime)
	if err := session.WritePbPacket(uint16(pb.OpCode_Op_S2C_Pong), &pb.S2C_Pong{
		OK: true,
	}); err != nil {
		return err
	}
	return nil
}

// C2S_LoginHandler 登录
func (m *GateWayM) C2S_LoginHandler(session *gameapi.Session, b []byte) error {
	var msg pb.C2S_Login
	proto.Unmarshal(b, &msg)

	log.Infof("Login: %s %s", msg.Account, msg.Password)
	if msg.Account == "admin" && msg.Password == "123456" {
		session.Valid()
		if err := session.WritePbPacket(uint16(pb.OpCode_Op_S2C_Login), &pb.S2C_Login{Ok: true}); err != nil {
			return err
		}
		return nil
	}

	session.Close()
	return errors.New("session is valid failed")
}
