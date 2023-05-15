package gameapi

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/trainking/goboot/pkg/log"
	"github.com/xtaci/kcp-go"

	"github.com/gorilla/websocket"
)

type (

	// NetConfig 配置
	NetConfig struct {
		Addr         string        // 监听地址
		TLSConfig    *tls.Config   // TLS配置
		WriteTimeout time.Duration // 写入超时
		ReadeTimeout time.Duration // 读取超时

		WebSocketPath string // ws连接的升级使用地址
		KcpMode       string // kcp的模式
	}

	// NetListener 接收请求的监听器
	NetListener interface {
		// Accept 接收一个连接
		Accept() (NetConn, error)

		// Close 关闭连接
		Close()
	}

	// NetConn 传输协议的抽象
	NetConn interface {

		// ReadPacket 读取报文
		ReadPacket() (Packet, error)

		// WritePacket 写入报文
		WritePacket(Packet) error

		// Close 关闭连接
		Close()
	}

	// TcpNetListener tcp连接监听器
	TcpNetListener struct {
		listener net.Listener
		config   *NetConfig // 配置
	}

	// KcpNetListener kcp连接监听器
	KcpNetListener struct {
		listener net.Listener
		config   *NetConfig
	}

	// WebSocketNetListener websocket连接的监听器
	WebSocketNetListener struct {
		server   *http.Server
		connChan chan *websocket.Conn
		config   *NetConfig

		ugrader websocket.Upgrader
	}

	// TcpNetConn tpc协议的连接抽象
	TcpNetConn struct {
		conn   net.Conn   // 连接
		config *NetConfig // 配置
	}

	// KcpNetConn kcp协议的连接抽象
	KcpNetConn struct {
		conn   net.Conn
		config *NetConfig // 配置
	}

	// WebSocketNetConn Websocket协议的连接抽象
	WebSocketNetConn struct {
		conn   *websocket.Conn
		config *NetConfig
	}
)

// NewTcpNetListener 创建一个tcp连接监听器
func NewTcpNetListener(config NetConfig) (NetListener, error) {
	l, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, err
	}

	return &TcpNetListener{listener: l, config: &config}, nil
}

// NewTcpNetConn 创建一个Tcp连接
func NewTcpNetConn(c net.Conn, config *NetConfig) NetConn {
	return &TcpNetConn{conn: c, config: config}
}

// Accept 接收连接
func (l *TcpNetListener) Accept() (NetConn, error) {
	c, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	if l.config.TLSConfig != nil {
		tlsConn := tls.Server(c, l.config.TLSConfig)

		return NewTcpNetConn(tlsConn, l.config), nil
	}

	return NewTcpNetConn(c, l.config), nil
}

// Close 关闭连接监听器
func (l *TcpNetListener) Close() {
	l.listener.Close()
}

// ReadPacket 读取数据报文
func (t *TcpNetConn) ReadPacket() (Packet, error) {
	if t.config.ReadeTimeout > 0 {
		t.conn.SetReadDeadline(time.Now().Add(t.config.ReadeTimeout))
	}

	return Packing(t.conn)
}

// WritePacket 写入数据报文
func (t *TcpNetConn) WritePacket(p Packet) error {
	if t.config.WriteTimeout > 0 {
		t.conn.SetWriteDeadline(time.Now().Add(t.config.WriteTimeout))
	}

	_, err := t.conn.Write(p.Serialize())
	return err
}

// Clost 关闭连接
func (t *TcpNetConn) Close() {
	t.conn.Close()
}

// NewKcpNetListener 创建一个kcp连接监听器
func NewKcpNetListener(config NetConfig) (NetListener, error) {
	l, err := kcp.Listen(config.Addr)
	if err != nil {
		return nil, err
	}
	return &KcpNetListener{listener: l, config: &config}, nil
}

// Accept 接收连接
func (l *KcpNetListener) Accept() (NetConn, error) {
	c, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	kcpConn := c.(*kcp.UDPSession)
	// 极速模式；普通模式参数为 0, 40, 0, 0
	if l.config.KcpMode == "nomarl" {
		kcpConn.SetNoDelay(0, 40, 0, 0)
	} else {
		kcpConn.SetNoDelay(1, 10, 2, 1)
	}

	kcpConn.SetStreamMode(true)
	kcpConn.SetWindowSize(4096, 4096)
	kcpConn.SetReadBuffer(4 * 65536 * 1024)
	kcpConn.SetWriteBuffer(4 * 65536 * 1024)
	kcpConn.SetACKNoDelay(true)

	if l.config.TLSConfig != nil {
		tlsConn := tls.Server(kcpConn, l.config.TLSConfig)
		return NewKcpNetConn(tlsConn, l.config), nil
	}

	return NewKcpNetConn(kcpConn, l.config), nil
}

// Close 关闭连接监听器
func (l *KcpNetListener) Close() {
	l.listener.Close()
}

// NewKcpNetConn 创建一个kcp连接
func NewKcpNetConn(c net.Conn, config *NetConfig) NetConn {
	return &KcpNetConn{conn: c, config: config}
}

// ReadPacket 读取数据包
func (k *KcpNetConn) ReadPacket() (Packet, error) {
	if k.config.ReadeTimeout > 0 {
		k.conn.SetReadDeadline(time.Now().Add(k.config.ReadeTimeout))
	}

	return Packing(k.conn)
}

// WritePacket 写入数据包
func (k *KcpNetConn) WritePacket(p Packet) error {
	if k.config.WriteTimeout > 0 {
		k.conn.SetWriteDeadline(time.Now().Add(k.config.WriteTimeout))
	}

	_, err := k.conn.Write(p.Serialize())
	return err
}

// Close 关闭连接
func (k *KcpNetConn) Close() {
	k.conn.Close()
}

// NewWebSocketNetListener 创建Websocket监听器
func NewWebSocketNetListener(config NetConfig) (NetListener, error) {
	l := &WebSocketNetListener{
		connChan: make(chan *websocket.Conn),
		config:   &config,
		ugrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	http.HandleFunc(config.WebSocketPath, l.handleWebSocket)

	server := &http.Server{
		Addr: config.Addr,
	}

	if config.TLSConfig != nil {
		server.TLSConfig = config.TLSConfig
		go func() {
			if err := server.ListenAndServeTLS("", ""); err != nil {
				log.Errorf("Websocket ListenAndServeTLS Eroor: %s", err)
			}
		}()

	} else {
		go func() {
			if err := server.ListenAndServe(); err != nil {
				log.Errorf("Websocket ListenAndServe Eroor: %s", err)
			}
		}()
	}

	l.server = server

	return l, nil
}

// handleWebSocket 将http升级成websocket
func (l *WebSocketNetListener) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := l.ugrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	l.connChan <- conn
}

// Accept 接收请求
func (l *WebSocketNetListener) Accept() (NetConn, error) {
	conn := <-l.connChan

	return NewWebSocketNetConn(conn, l.config), nil
}

// Close 关闭连接
func (l *WebSocketNetListener) Close() {
	close(l.connChan)
	l.server.Close()
}

// NewWebSocketNetConn 读取WebSocket连接
func NewWebSocketNetConn(c *websocket.Conn, config *NetConfig) NetConn {
	return &WebSocketNetConn{conn: c, config: config}
}

// ReadPacket 读取数据包
func (w *WebSocketNetConn) ReadPacket() (Packet, error) {
	if w.config.ReadeTimeout > 0 {
		w.conn.SetReadDeadline(time.Now().Add(w.config.ReadeTimeout))
	}

	_, message, err := w.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return NewPacket(message), nil
}

// WritePacket 写入数据包
func (w *WebSocketNetConn) WritePacket(p Packet) error {
	if w.config.WriteTimeout > 0 {
		w.conn.SetWriteDeadline(time.Now().Add(w.config.WriteTimeout))
	}

	return w.conn.WriteMessage(websocket.BinaryMessage, p.Serialize())
}

// Close 关闭连接
func (w *WebSocketNetConn) Close() {
	if w.conn != nil {
		w.conn.Close()
	}
}
