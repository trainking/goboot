package gameapi

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

// Packet 包接口
type Packet interface {
	// Serialize 序列化
	Serialize() []byte

	// OpeCode 获取该包的OpCode
	OpCode() uint16

	// BodyLen 内容长度
	BodyLen() uint16

	// Body 获取完整body
	Body() []byte
}

// PbPacket 基于Protobuffer的包协议
type PbPacket struct {
	buff []byte
}

// NewPbPacket 新建一个pb的Packet
func NewPbPacket(buff []byte, opcode uint16) *PbPacket {
	p := new(PbPacket)

	p.buff = make([]byte, 4+len(buff))
	binary.BigEndian.PutUint16(p.buff[0:2], uint16(len(buff)))
	binary.BigEndian.PutUint16(p.buff[2:4], opcode)
	copy(p.buff[4:], buff)

	return p
}

// Serialize 序列化，输出完整的字符数组
func (p *PbPacket) Serialize() []byte {
	return p.buff
}

// OpCode 包的2-3位为OpCode
func (p *PbPacket) OpCode() uint16 {
	return binary.BigEndian.Uint16(p.buff[2:4])
}

// BodyLen 报文内容长度
func (p *PbPacket) BodyLen() uint16 {
	return binary.BigEndian.Uint16(p.buff[0:2])
}

// Body 读取body所有字符
func (p *PbPacket) Body() []byte {
	return p.buff[4:]
}

// ReadPacket 读取数据打包
func ReadPacket(r io.Reader) (Packet, error) {
	// 4字节头
	var headrBytes = make([]byte, 4)

	// 读取头
	if _, err := io.ReadFull(r, headrBytes); err != nil {
		return nil, err
	}

	opcode := binary.BigEndian.Uint16(headrBytes[2:4])

	bodyLength := binary.BigEndian.Uint16(headrBytes[0:2])
	buff := make([]byte, bodyLength)

	// 读取body
	if _, err := io.ReadFull(r, buff); err != nil {
		return nil, err
	}

	return NewPbPacket(buff, opcode), nil
}

// CretaePbPacket 创建要给protobuf的包
func CretaePbPacket(opcode uint16, msg proto.Message) (Packet, error) {
	msgB, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return NewPbPacket(msgB, opcode), nil
}
