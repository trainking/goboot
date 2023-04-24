package gameapi

import (
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"
)

// HeartPacket 默认使用0号协议作为心跳包的协议
var HeartPacket Packet = NewDefaultPacket(nil, 0)

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

// DefaultPacket 基于Protobuffer的包协议
type DefaultPacket struct {
	buff []byte
}

// NewPbPacket 新建一个pb的Packet
func NewDefaultPacket(buff []byte, opcode uint16) *DefaultPacket {
	p := new(DefaultPacket)

	p.buff = make([]byte, 4+len(buff))
	binary.BigEndian.PutUint16(p.buff[0:2], uint16(len(buff)))
	binary.BigEndian.PutUint16(p.buff[2:4], opcode)
	if len(buff) > 0 {
		copy(p.buff[4:], buff)
	}

	return p
}

// NewPacket 从一份数据字节中构造Packet
func NewPacket(buff []byte) Packet {
	p := new(DefaultPacket)
	p.buff = buff

	return p
}

// Serialize 序列化，输出完整的字符数组
func (p *DefaultPacket) Serialize() []byte {
	return p.buff
}

// OpCode 包的2-3位为OpCode
func (p *DefaultPacket) OpCode() uint16 {
	return binary.BigEndian.Uint16(p.buff[2:4])
}

// BodyLen 报文内容长度
func (p *DefaultPacket) BodyLen() uint16 {
	return binary.BigEndian.Uint16(p.buff[0:2])
}

// Body 读取body所有字符
func (p *DefaultPacket) Body() []byte {
	return p.buff[4:]
}

// Packing 从io流中读取出包
func Packing(r io.Reader) (Packet, error) {
	// 4字节头
	var headrBytes = make([]byte, 4)

	// 读取头
	if _, err := io.ReadFull(r, headrBytes); err != nil {
		return nil, err
	}

	opcode := binary.BigEndian.Uint16(headrBytes[2:4])
	bodyLength := binary.BigEndian.Uint16(headrBytes[0:2])

	var buff []byte
	if bodyLength > 0 {
		buff = make([]byte, bodyLength)

		// 读取body
		if _, err := io.ReadFull(r, buff); err != nil {
			return nil, err
		}
	}

	return NewDefaultPacket(buff, opcode), nil
}

// CretaePbPacket 创建要给protobuf的包
func CretaePbPacket(opcode uint16, msg proto.Message) (Packet, error) {
	msgB, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return NewDefaultPacket(msgB, opcode), nil
}
