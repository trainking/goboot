package gateway

import (
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

}

func (m *GateWayM) Group() map[uint16]gameapi.Handler {
	return map[uint16]gameapi.Handler{
		uint16(pb.OpCode_Ping): m.C2S_PingHandler,
	}
}

// C2S_PingHandler
func (m *GateWayM) C2S_PingHandler(session *gameapi.Session, b []byte) error {

	var msg pb.C2S_Ping
	proto.Unmarshal(b, &msg)

	log.Infof("Receive %v", msg.TickTime)
	if err := session.WritePbPacket(uint16(pb.OpCode_Pong), &pb.S2C_Pong{
		OK: true,
	}); err != nil {
		return err
	}
	return nil
}
