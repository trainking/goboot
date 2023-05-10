package gameapi

import (
	"context"
	"crypto/tls"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/etcdx"
	"github.com/trainking/goboot/pkg/log"
	"github.com/trainking/goboot/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type (
	// App 游戏服务器实现
	App struct {
		boot.BaseInstance

		gd             GameMetaData          // 服务器元数据
		un             *UserNats             // 用户消息分发
		listener       NetListener           // 网络监听
		exitChan       chan struct{}         // 退出信号
		waitGroup      sync.WaitGroup        // 等待携程控制
		closeOnce      sync.Once             // 保证关闭只执行一次
		serviceManager *etcdx.ServiceManager // etcd注册管理器

		connectListener    Listener     // 会话建立时的连接监听
		disconnectListener Listener     // 会话断开时的连接监听器
		beforeMiddleware   []Middleware // 会话处理消息之前执行的中间件
		afterMiddleware    []Middleware // 会话处理消息之后执行的中间件

		router map[uint16]Handler // 处理器映射

		totalConn int64 // 连接数量

		modules []Moddule // 服务上的模块

		sessions  map[int64]*Session // 有效Session映射
		sessionMu sync.RWMutex
	}

	// Middleware 中间件处理
	Middleware struct {
		// Condition 是否要处理的opcode
		Condition func(uint16) bool
		// Do 处理执行
		Do func(Context) error
	}

	// Listener 监听器
	Listener func(*Session) error

	// Module 模块
	Moddule interface {

		// 初始化模块
		Init(app *App)
	}

	// Handler 处理类型
	Handler func(Context) error

	// SessionCreator session创建生成器，做一些预处理
	SessionCreator func(net.Conn, *App) *Session
)

// New 创建一个游戏服务器接口实例
func New(name string, configPath string, addr string, instancdID int64) *App {
	// 加载配置
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	un, err := NewUserNats(v.GetString("NatsUrl"), name)
	if err != nil {
		panic(err)
	}

	app := new(App)
	app.Name = name
	app.gd = GameMetaData{
		ID:       instancdID,
		State:    StateZero,
		Password: v.GetString("Password"),
	}
	app.un = un
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
func (a *App) AddBeforeMiddleware(m Middleware) {
	a.beforeMiddleware = append(a.beforeMiddleware, m)
}

// AddAfterMiddleware 设置详细处理后的中间件
func (a *App) AddAfterMiddleware(m Middleware) {
	a.afterMiddleware = append(a.afterMiddleware, m)
}

// AddHandler 增加处理器
func (a *App) AddHandler(opcode interface{}, h Handler) {
	_op := opcodeChange(opcode)
	if _op == 0 {
		return
	}
	a.router[_op] = h
}

// AddModule 增加所有模块
func (a *App) AddModule(module Moddule) {
	a.modules = append(a.modules, module)
}

// AddSession 加入Session
func (a *App) AddSession(userID int64, session *Session) {
	a.sessionMu.Lock()
	defer a.sessionMu.Unlock()

	// 已有Session，关闭
	if oldSession, ok := a.sessions[userID]; ok {
		oldSession.Close()
	}
	a.sessions[userID] = session
}

// DelSession 删除Session
func (a *App) DelSession(userID int64) {
	a.sessionMu.Lock()
	defer a.sessionMu.Unlock()

	delete(a.sessions, userID)
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
	err = a.sendActorLocation(userID, p)
	if err == ErrUserNoIn {
		// 发送给其他实例
		return a.sendActorPush(userID, p)
	}

	return err
}

// SendAllActor 向所有玩家发送消息
func (a *App) SendAllActor(opcode uint16, msg proto.Message) error {
	p, err := CretaePbPacket(opcode, msg)
	if err != nil {
		return err
	}

	return a.sendActorPush(-1, p)
}

// sendActorLocation 发送消息给本地实例
func (a *App) sendActorLocation(userID int64, p Packet) error {
	a.sessionMu.RLock()
	defer a.sessionMu.RUnlock()
	if session, ok := a.sessions[userID]; ok {
		if !session.IsClosed() {
			// 如果定义handler，必须发送给handler处理
			if _, ok := a.router[p.OpCode()]; ok {
				session.receiveChan <- p
				return nil
			}
			return session.WritePacket(p)
		}
	}

	return ErrUserNoIn
}

// sendAllActorLocation 像所有的本地玩家发送消息
func (a *App) sendAllActorLocation(m PushActorMessage) error {
	a.sessionMu.RLock()
	defer a.sessionMu.RUnlock()
	for userID := range a.sessions {
		session := a.sessions[userID]
		if !session.IsClosed() {
			p := NewDefaultPacket(m.Msg, m.OpCode)
			// 如果定义handler，必须发送给handler处理
			if _, ok := a.router[p.OpCode()]; ok {
				session.receiveChan <- p
			} else {
				if err := session.WritePacket(p); err != nil {
					return err
				}
			}
		}
	}

	return ErrUserNoIn
}

// sendActorPush 发送消息给远程实例
func (a *App) sendActorPush(userID int64, p Packet) error {
	return a.un.Publish(userID, p)
}

// Init 初始化服务
func (a *App) Init() (err error) {
	if err = a.BaseInstance.Init(); err != nil {
		return err
	}

	// 初始化各个模块
	for _, m := range a.modules {
		m.Init(a)
	}

	var netConfig = NetConfig{
		Addr:         a.Addr,
		WriteTimeout: time.Duration(a.Config.GetInt("ConnWriteTimeout") * int(time.Second)),
		ReadeTimeout: time.Duration(a.Config.GetInt("ConnReadTimeout") * int(time.Second)),
	}

	// 加载tls配置
	tlsConfigMap := a.Config.GetStringMapString("TLS")
	if len(tlsConfigMap) > 0 {
		cert, err := tls.LoadX509KeyPair(tlsConfigMap["certfile"], tlsConfigMap["keyfile"])
		if err != nil {
			return err
		}
		netConfig.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		a.gd.UseTLS = true
	}

	// 根据配置传输层协议
	network := a.Config.GetString("Network")
	switch network {
	case "tcp":
		a.listener, err = NewTcpNetListener(netConfig)
		if err != nil {
			return err
		}
	case "kcp":
		a.listener, err = NewKcpNetListener(netConfig)
		if err != nil {
			return err
		}
	case "websocket":
		netConfig.WebSocketPath = a.Config.GetString("WebsocketPath")
		a.listener, err = NewWebSocketNetListener(netConfig)
		if err != nil {
			return err
		}
	default:
		return errors.Wrap(ErrNoImplementNetwork, network)
	}
	a.gd.Network = network

	// 监听广播的其他实例消息
	go a.subscribePushUserMsg()

	// 注册服务
	if err := a.registerEtcd(); err != nil {
		return err
	}

	return nil
}

// subscribePushUserMsg 订阅消费推送玩家的消息
func (a *App) subscribePushUserMsg() {
	// 开启消费, 消费limit使用接收消息的一半
	if err := a.un.StartSubscribe(a.Config.GetInt("ReceiveLimit") / 2); err != nil {
		log.Errorf("StartSubscribe Error: %v", err)
		return
	}
	defer func() {
		a.un.Close()
		a.Stop()
	}()

	for {
		select {
		case <-a.exitChan:
			return
		case m := <-a.un.Consume():
			if m.UserID != -1 {
				p := NewDefaultPacket(m.Msg, m.OpCode)
				// 向特定玩家发送消息
				if err := a.sendActorLocation(m.UserID, p); err != nil && err != ErrUserNoIn {
					log.Errorf("PushActorMessage sendActorLocation Error: %v", err)
					return
				}
			} else {
				// 向所有玩家广播消息
				if err := a.sendAllActorLocation(m); err != nil && err != ErrUserNoIn {
					log.Errorf("PushActorMessage sendAllActorLocation Error: %v", err)
					return
				}
			}

		}
	}
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
			NewSession(conn, a).Run()
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
		a.serviceManager.Destory(a.Addr)
	})

	// 等待所有携程执行完
	a.waitGroup.Wait()
}

// OnConnect 连接时处理
func (a *App) OnConnect(session *Session) bool {
	atomic.AddInt64(&a.totalConn, 1)
	session.startValidTimer()

	if a.connectListener != nil {
		if err := a.connectListener(session); err != nil {
			return false
		}
	}

	return true
}

// OnMessage 消息处理
func (a *App) OnMessage(session *Session, p Packet) bool {
	// 消息的分发
	if h, ok := a.router[p.OpCode()]; ok {
		go func() {
			// 释放Packet
			defer p.Free()

			ctx := NewDefaultContext(context.Background(), a, session, p.OpCode(), p.Body())
			// 处理消息之前，中间件过滤
			for _, m := range a.beforeMiddleware {
				if !m.Condition(p.OpCode()) {
					continue
				}

				if err := m.Do(ctx); err != nil {
					log.Errorf("Middle %d is Error: %v", p.OpCode(), err)
					continue
				}
			}

			// 处理消息
			if err := h(ctx); err != nil {
				log.Errorf("Handler %d Error: %s ", p.OpCode(), err)
				return
			}

			// 处理之后，中间件操作
			for _, m := range a.afterMiddleware {
				if !m.Condition(p.OpCode()) {
					continue
				}

				if err := m.Do(ctx); err != nil {
					log.Errorf("Middle %d is Error: %v", p.OpCode(), err)
					continue
				}
			}
		}()
	} else {
		log.Errorf("Handler %d is No Handler", p.OpCode())
		return false
	}
	return true
}

// OnDisConnect 断线处理
func (a *App) OnDisConnect(sesssion *Session) {
	atomic.AddInt64(&a.totalConn, -1)

	if a.disconnectListener != nil {
		a.disconnectListener(sesssion)
	}

	// 从有效连接中删除
	if sesssion.IsValid() {
		a.DelSession(sesssion.UserID())
	}
}

// opcodeChange opcode的类型转换
func opcodeChange(opcode interface{}) uint16 {
	var _op uint16
	switch reflect.TypeOf(opcode).Kind() {
	case reflect.Int32, reflect.Int, reflect.Int64:
		_op = uint16(reflect.ValueOf(opcode).Int())
	case reflect.Uint, reflect.Uint16, reflect.Uint32:
		_op = uint16(reflect.ValueOf(opcode).Uint())
	default:
		log.Errorf("wrong opcode %v kind: %v", opcode, reflect.TypeOf(opcode).Kind())
	}
	return _op
}
