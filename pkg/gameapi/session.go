package gameapi

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/trainking/goboot/pkg/idgen"
	"google.golang.org/protobuf/proto"
)

// Session 网络会话
type Session struct {
	ID  int64 // session的唯一标识
	app *App

	conn       NetConn       // 网络连接
	closeOnce  sync.Once     // 关闭控制
	closeFlag  int32         // 关闭标志
	closeChan  chan struct{} // 关闭channel
	isValid    int32         // 是否验证为有效连接
	validTimer *time.Timer   // 设置验证超时
	expired    int64         // Session过期时间

	sendChan    chan Packet // 发送队列
	receiveChan chan Packet // 接收队列

	callback SessionCallback // 回调函数

	userID   int64       // sesssion所属的玩家ID
	userData interface{} // session所属玩家额外数据

	heartbeatCount int64 // 连接的心跳包计数器
}

// SessionCallback session触发外部事件调用
type SessionCallback interface {

	// OnConnect 当连接建立时调用
	OnConnect(*Session) bool

	// OnHeartbeat 当心跳包到达时处理
	OnHeartbeat(*Session)

	// OnMessage 当连接处理消息时
	OnMessage(*Session, Packet) bool

	// OnDisConnect 当连接断开时
	OnDisConnect(*Session)
}

// NewSession 新建一个Session
func NewSession(conn NetConn, app *App) *Session {
	return &Session{
		ID:          idgen.GenerateInt64(),
		app:         app,
		callback:    app,
		conn:        conn,
		closeChan:   make(chan struct{}),
		sendChan:    make(chan Packet, app.Config.GetInt("SendLimit")),
		receiveChan: make(chan Packet, app.Config.GetInt("ReceiveLimit")),
	}
}

// Close 关闭Session
func (s *Session) Close() {
	s.closeOnce.Do(func() {
		atomic.StoreInt32(&s.closeFlag, 1)
		close(s.closeChan)
		close(s.sendChan)
		close(s.receiveChan)
		s.conn.Close()
		go s.callback.OnDisConnect(s)
	})
}

// IsClosed 是否关闭
func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closeFlag) > 0
}

// startValidTimer 开始验证计时
func (s *Session) startValidTimer() {
	s.validTimer = time.NewTimer(time.Second * time.Duration(s.app.Config.GetInt("ValidTimeout")))
	go func() {
		select {
		// 关服退出
		case <-s.app.exitChan:
			return
		case <-s.closeChan:
			return
		case <-s.validTimer.C:
			if !s.validTimer.Stop() && !s.IsValid() {
				s.Close()
			}
			return
		}
	}()
}

func (s *Session) App() *App {
	return s.app
}

// UserID 获取玩家的UserID
func (s *Session) UserID() int64 {
	return s.userID
}

// UserData 获取玩家的额外数据
func (s *Session) UserData() interface{} {
	return s.userData
}

// HeartbeatCount 连接的心跳包计数
func (s *Session) HeartbeatCount() int64 {
	return atomic.LoadInt64(&s.heartbeatCount)
}

// SetUserData 设置玩家的额外数据
func (s *Session) SetUserData(d interface{}) {
	s.userData = d
}

// SetUserID 设置玩家ID
func (s *Session) SetUserID(userID int64) {
	s.userID = userID
}

// SetExpired 设置session的过期时间
func (s *Session) SetExpired(expired int64) {
	s.expired = expired
}

// Send 向此session推送消息
func (s *Session) Send(opcode interface{}, msg proto.Message) error {
	_op := opcodeChange(opcode)
	if _op == 0 {
		return ErrWrongOpCode
	}
	p, err := CretaePbPacket(_op, msg)
	if err != nil {
		return err
	}

	return s.WritePacket(p)
}

// IsExipred 判断Session是否已经过期
func (s *Session) IsExipred() bool {
	if s.expired > 0 {
		if time.Now().Unix() > s.expired {
			return true
		}
	}

	return false
}

// IsValid 是否验证为有效连接
func (s *Session) IsValid() bool {
	return atomic.LoadInt32(&s.isValid) > 0
}

// Valid 设置为有效连接
func (s *Session) valid() {
	if s.validTimer != nil {
		s.validTimer.Stop()
	}

	atomic.StoreInt32(&s.isValid, 1)
}

// WritePacket 写入发送包
func (s *Session) WritePacket(p Packet) (err error) {
	if s.IsClosed() {
		return ErrConnClosing
	}

	defer func() {
		if e := recover(); e != nil {
			err = ErrConnClosing
		}
	}()

	select {
	case s.sendChan <- p:
		return nil
	case <-s.closeChan:
		return ErrConnClosing
	}
}

// Run 执行主体逻辑，三套循环
// readLoop 读循环
// writeLoop 写循环
// handleLopp 处理消息循环
func (s *Session) Run() {
	if !s.callback.OnConnect(s) {
		return
	}

	asyncDo(s.handleLoop, &s.app.waitGroup)
	asyncDo(s.readLoop, &s.app.waitGroup)
	asyncDo(s.writeLoop, &s.app.waitGroup)
}

// readLoop 读循环
func (s *Session) readLoop() {
	defer func() {
		recover()
		s.Close()
	}()

	for {
		select {
		// 关服退出
		case <-s.app.exitChan:
			return
		case <-s.closeChan:
			return
		default:
		}

		if s.IsExipred() {
			return
		}

		p, err := s.conn.ReadPacket()
		if err != nil {
			return
		}

		s.receiveChan <- p
	}
}

// writeLoop 写循环
func (s *Session) writeLoop() {
	defer func() {
		recover()
		s.Close()
	}()

	for {
		select {
		// 关服退出
		case <-s.app.exitChan:
			return
		case <-s.closeChan:
			return
		case p := <-s.sendChan:
			if s.IsClosed() {
				return
			}

			if err := s.conn.WritePacket(p); err != nil {
				return
			}
			p.Free()
		}
	}
}

// handleLoop 处理循环
func (s *Session) handleLoop() {
	defer func() {
		recover()
		s.Close()
	}()

	for {
		select {
		// 关服退出
		case <-s.app.exitChan:
			return
		case <-s.closeChan:
			return
		case p := <-s.receiveChan:
			if s.IsClosed() {
				return
			}

			// 处理心跳包
			if p.OpCode() == 0 {
				// 心跳包过载，则跳出
				s.callback.OnHeartbeat(s)
				continue
			}

			if !s.callback.OnMessage(s, p) {
				return
			}
		}
	}
}

// asyncDo 异步执行循环
func asyncDo(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}
