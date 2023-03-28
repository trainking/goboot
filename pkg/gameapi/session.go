package gameapi

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Session 网络会话
type Session struct {
	app *App

	conn      net.Conn      // 网络连接
	closeOnce sync.Once     // 关闭控制
	closeFlag int32         // 关闭标志
	closeChan chan struct{} // 关闭channel

	sendChan    chan Packet // 发送队列
	receiveChan chan Packet // 接收队列

	callback SessionCallback // 回调函数
}

// SessionCallback session触发外部事件调用
type SessionCallback interface {

	// OnConnect 当连接建立时调用
	OnConnect(*Session) bool

	// OnMessage 当连接处理消息时
	OnMessage(*Session, Packet) bool

	// OnDisConnect 当连接断开时
	OnDisConnect(*Session)
}

// NewSession 新建一个Session
func NewSession(conn net.Conn, app *App) *Session {
	return &Session{
		app:         app,
		callback:    app.callback,
		conn:        conn,
		closeChan:   make(chan struct{}),
		sendChan:    make(chan Packet, app.Config.GetInt("SendLimit")),
		receiveChan: make(chan Packet, app.Config.GetInt("ReceiveLimit")),
	}
}

// GetConn 获取连接实例
func (s *Session) GetConn() net.Conn {
	return s.conn
}

// Close 关闭Session
func (s *Session) Close() {
	s.closeOnce.Do(func() {
		atomic.StoreInt32(&s.closeFlag, 1)
		close(s.closeChan)
		close(s.sendChan)
		close(s.receiveChan)
		s.conn.Close()
		s.callback.OnDisConnect(s)
	})
}

// IsClosed 是否关闭
func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closeFlag) > 0
}

// Run 执行主体逻辑，三套循环
// readLoop 读循环
// writeLoop 写循环
// handleLopp 处理消息循环
func (s *Session) Run() {
	if !s.callback.OnConnect(s) {
		return
	}

	asyncDo(s.handleLoop, s.app.waitGroup)
	asyncDo(s.readLoop, s.app.waitGroup)
	asyncDo(s.writeLoop, s.app.waitGroup)
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

		// 设置连接读取超时
		s.conn.SetReadDeadline(time.Now().Add(time.Duration(s.app.Config.GetInt("ConnReadTimeout")) * time.Second))
		p, err := ReadPacket(s.conn)
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
			s.conn.SetWriteDeadline(time.Now().Add(time.Duration(s.app.Config.GetInt("ConnWriteTimeout")) * time.Second))
			if _, err := s.conn.Write(p.Serialize()); err != nil {
				return
			}
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
