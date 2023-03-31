package gameapi

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/log"
	"github.com/trainking/goboot/pkg/utils"
	"github.com/xtaci/kcp-go"
)

type (
	// App 游戏服务器实现
	App struct {
		boot.BaseInstance

		listener  net.Listener   // 网络监听
		exitChan  chan struct{}  // 退出信号
		waitGroup sync.WaitGroup // 等待携程控制
		closeOnce sync.Once      // 保证关闭只执行一次

		connectListener   Listener       // 会话建立时的连接监听器
		disconnectListner Listener       // 会话断开时的连接监听器
		creator           SessionCreator // 会话创建时，预处理的生成器

		router map[uint16]Handler // 处理器映射

		totalConn int64 // 连接数量
	}

	// Listener 监听器
	Listener interface {
		Do(*Session) error
	}

	// Handler 处理类型
	Handler func(*Session, []byte) error

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
	return app
}

// SetConnectListener 设置连接监听器
func (a *App) SetConnectListener(l Listener) {
	a.connectListener = l
}

// SetDisconnectListener 设置断连监听器
func (a *App) SetDisconnectListener(l Listener) {
	a.disconnectListner = l
}

// AddHandler 增加处理器
func (a *App) AddHandler(opcode uint16, h Handler) {
	a.router[opcode] = h
}

// GetTotalConn 获取连接总数
func (a *App) GetTotalConn() int64 {
	return atomic.LoadInt64(&a.totalConn)
}

// Init 初始化服务
func (a *App) Init() (err error) {
	if err = a.BaseInstance.Init(); err != nil {
		return err
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
		if err := a.connectListener.Do(session); err != nil {
			return false
		}
	}

	atomic.AddInt64(&a.totalConn, 1)
	return true
}

// OnMessage 消息处理
func (a *App) OnMessage(session *Session, p Packet) bool {
	// 消息的分发
	if h, ok := a.router[p.OpCode()]; ok {
		go func() {
			if err := h(session, p.Body()); err != nil {
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

	if a.disconnectListner != nil {
		a.disconnectListner.Do(sesssion)
	}
}
