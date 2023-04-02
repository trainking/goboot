package gameapi

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/log"
	"github.com/trainking/goboot/pkg/utils"
	"github.com/xtaci/kcp-go"
	"google.golang.org/protobuf/proto"
)

type (
	// App 游戏服务器实现
	App struct {
		boot.BaseInstance

		listener  net.Listener   // 网络监听
		exitChan  chan struct{}  // 退出信号
		waitGroup sync.WaitGroup // 等待携程控制
		closeOnce sync.Once      // 保证关闭只执行一次

		connectListener    Listener       // 会话建立时的连接监听器
		disconnectListener Listener       // 会话断开时的连接监听器
		beforeMiddleware   Middleware     // 会话处理消息之前执行的监听器
		creator            SessionCreator // 会话创建时，预处理的生成器

		router map[uint16]Handler // 处理器映射

		totalConn int64 // 连接数量

		modules []Moddule // 服务上的模块

		sessions map[int64]*Session // 有效Session映射
	}

	// Middleware 中间件处理
	Middleware func(*Session, Packet) error

	// Listener 监听器
	Listener func(*Session) error

	// Module 模块
	Moddule interface {

		// 初始化模块
		Init(app *App)

		// 模块的分组路由
		Group() map[uint16]Handler
	}

	// Handler 处理类型
	Handler func(Context) error

	// SessionCreator session创建生成器，做一些预处理
	SessionCreator func(net.Conn, *App) *Session
)

// New 创建一个游戏服务器接口实例
func New(configPath string, addr string, instancdID int64) *App {
	// 加载配置
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	app := new(App)
	app.Config = v
	app.Addr = addr
	app.IntanceID = instancdID
	app.router = make(map[uint16]Handler)
	app.exitChan = make(chan struct{})
	app.sessions = make(map[int64]*Session)
	return app
}

// SetConnectListener 设置连接监听器
func (a *App) SetConnectListener(l Listener) {
	a.connectListener = l
}

// SetDisconnectListener 设置断连监听器
func (a *App) SetDisconnectListener(l Listener) {
	a.disconnectListener = l
}

// SetBeforeMiddleware 设置消息处理前中间件
func (a *App) SetBeforeMiddleware(m Middleware) {
	a.beforeMiddleware = m
}

// AddHandler 增加处理器
func (a *App) AddHandler(opcode uint16, h Handler) {
	a.router[opcode] = h
}

// AddModule 增加所有模块
func (a *App) AddModule(module Moddule) {
	a.modules = append(a.modules, module)
}

// GetTotalConn 获取连接总数
func (a *App) GetTotalConn() int64 {
	return atomic.LoadInt64(&a.totalConn)
}

// SendActor 向指定玩家发送消息
func (a *App) SendActor(userID int64, opcode uint16, msg proto.Message) error {
	p, err := CretaePbPacket(opcode, msg)
	if err != nil {
		return err
	}

	// 如果在本实例，则直接发送
	if session, ok := a.sessions[userID]; ok {
		return session.WritePacket(p)
	}

	// TODO：发送给其他实例
	return nil
}

// Init 初始化服务
func (a *App) Init() (err error) {
	if err = a.BaseInstance.Init(); err != nil {
		return err
	}

	// 初始化各个模块
	for _, m := range a.modules {
		m.Init(a)

		// 加入路由映射
		for k, v := range m.Group() {
			a.AddHandler(k, v)
		}
	}

	// 根据配置传输层协议
	network := a.Config.GetString("Network")
	switch network {
	case "tcp":
		a.listener, err = net.Listen("tcp", a.Addr)
		if err != nil {
			return err
		}
	case "kcp":
		a.listener, err = kcp.Listen(a.Addr)
		if err != nil {
			return err
		}

		// 设置kcp连接参数
		a.creator = func(c net.Conn, a *App) *Session {
			kcpConn := c.(*kcp.UDPSession)
			// 极速模式；普通模式参数为 0, 40, 0, 0
			kcpConn.SetNoDelay(1, 10, 2, 1)
			kcpConn.SetStreamMode(true)
			kcpConn.SetWindowSize(4096, 4096)
			kcpConn.SetReadBuffer(4 * 65536 * 1024)
			kcpConn.SetWriteBuffer(4 * 65536 * 1024)
			kcpConn.SetACKNoDelay(true)

			return NewSession(c, a)
		}
	default:
		return errors.Wrap(ErrNoImplementNetwork, network)
	}

	return nil
}

// Start 启动服务
func (a *App) Start() error {
	a.waitGroup.Add(1)
	defer func() {
		a.waitGroup.Done()
	}()

	for {
		select {
		case <-a.exitChan:
			return nil
		default:
		}

		conn, err := a.listener.Accept()
		if err != nil {
			continue
		}

		a.waitGroup.Add(1)
		go func() {
			// 处理连接数据
			if a.creator != nil {
				a.creator(conn, a).Run()
			} else {
				NewSession(conn, a).Run()
			}
			a.waitGroup.Done()
		}()
	}
}

// Stop 停止服务
func (a *App) Stop() {
	// 关闭资源
	a.closeOnce.Do(func() {
		close(a.exitChan)
		a.listener.Close()
	})

	// 等待所有携程执行完
	a.waitGroup.Wait()
}

// OnConnect 连接时处理
func (a *App) OnConnect(session *Session) bool {
	if a.connectListener != nil {
		if err := a.connectListener(session); err != nil {
			return false
		}
	}

	session.startValidTimer()
	atomic.AddInt64(&a.totalConn, 1)
	return true
}

// OnMessage 消息处理
func (a *App) OnMessage(session *Session, p Packet) bool {
	// 处理消息之前的预处理
	if a.beforeMiddleware != nil {
		if err := a.beforeMiddleware(session, p); err != nil {
			log.Errorf("beforeMiddleware Error: %v", err)
			return false
		}
	}
	// 消息的分发
	if h, ok := a.router[p.OpCode()]; ok {
		go func() {
			if err := h(NewDefaultContext(context.Background(), a, session, p.Body())); err != nil {
				log.Errorf("Handler %d Error: %s ", p.OpCode(), err)
			}
		}()
	} else {
		log.Errorf("Handler %d is No Handler", p.OpCode())
		return false
	}
	return true
}

// OnDisConnect 短线处理
func (a *App) OnDisConnect(sesssion *Session) {
	atomic.AddInt64(&a.totalConn, -1)

	if a.disconnectListener != nil {
		a.disconnectListener(sesssion)
	}
}
