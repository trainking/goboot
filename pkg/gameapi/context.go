package gameapi

import (
	"context"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type (
	// Context 抽象每个Handler的调用参数
	Context interface {

		// Context 返回一个context.Context
		Context() context.Context

		// Parmas 取出请求参数
		Params(m protoreflect.ProtoMessage) error

		// Session 获取这个玩家的Session
		Session() *Session

		// Send 发送消息到玩家
		Send(opcode uint16, msg proto.Message) error
	}

	// DefaultContext 默认Context实现
	DefaultContext struct {
		ctx     context.Context
		body    []byte
		session *Session
	}
)

func NewDefaultContext(ctx context.Context, body []byte, session *Session) Context {
	return &DefaultContext{
		ctx:     ctx,
		body:    body,
		session: session,
	}
}

func (c *DefaultContext) Context() context.Context {
	return c.ctx
}

func (c *DefaultContext) Params(m protoreflect.ProtoMessage) error {
	return proto.Unmarshal(c.body, m)
}

func (c *DefaultContext) Session() *Session {
	return c.session
}

// WritePbPacket 写入Protobuf的包
func (c *DefaultContext) Send(opcode uint16, msg proto.Message) error {
	p, err := CretaePbPacket(opcode, msg)
	if err != nil {
		return err
	}

	return c.session.WritePacket(p)
}
