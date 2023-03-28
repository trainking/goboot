package gameapi

import "net"

type Packet struct {
}

func (p *Packet) Serialize() []byte {
	return []byte{}
}

func ReadPacket(conn net.Conn) (Packet, error) {
	return Packet{}, nil
}
