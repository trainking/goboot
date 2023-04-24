package gameapi

import (
	"context"

	"github.com/google/uuid"
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
		Send(opcode interface{}, msg proto.Message) error

		// SendActor 向其他玩家发送消息
		SendActor(userID int64, opcode interface{}, msg proto.Message) error

		// Valid 验证玩家成功，传入用户ID
		Valid(userID int64)

		// GetOpCode 获取此次请求的opcode
		GetOpCode() uint16

		// GetRequestID 获得请求ID
		GetRequestID() string
	}

	// DefaultContext 默认Context实现
	DefaultContext struct {
		a         *App
		session   *Session
		ctx       context.Context
		opcode    uint16
		body      []byte
		requestID string
	}
)

// NewDefaultContext
func NewDefaultContext(ctx context.Context, a *App, session *Session, opcode uint16, body []byte) Context {
	id := uuid.New()
	return &DefaultContext{
		a:         a,
		session:   session,
		ctx:       ctx,
		opcode:    opcode,
		body:      body,
		requestID: id.String(),
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
func (c *DefaultContext) Send(opcode interface{}, msg proto.Message) error {
	_op := opcodeChange(opcode)
	if _op == 0 {
		return ErrWrongOpCode
	}
	p, err := CretaePbPacket(_op, msg)
	if err != nil {
		return err
	}

	return c.session.WritePacket(p)
}

// SendActor 向指定玩家发送消息
func (c *DefaultContext) SendActor(userID int64, opcode interface{}, msg proto.Message) error {
	_op := opcodeChange(opcode)
	if _op == 0 {
		return ErrWrongOpCode
	}
	return c.a.SendActor(userID, _op, msg)
}

// Valid 验证成功
func (c *DefaultContext) Valid(userID int64) {
	c.a.AddSession(userID, c.session)
	c.session.valid()
	c.session.SetUserID(userID)
}

// GetOpCode 获取此次处理的opcode
func (c *DefaultContext) GetOpCode() uint16 {
	return c.opcode
}

// GetRequestID 获得请求ID
func (c *DefaultContext) GetRequestID() string {
	return c.requestID
}

// GetRequestID 别名函数，保持一致
func GetRequestID(ctx Context) string {
	return ctx.GetRequestID()
}
