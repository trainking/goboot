package gameapi

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/xtaci/kcp-go"
)

type (

	// NetConfig 配置
	NetConfig struct {
		Addr         string        // 监听地址
		TLSConfig    *tls.Config   // TLS配置
		WriteTimeout time.Duration // 写入超时
		ReadeTimeout time.Duration // 读取超时
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

	// TcpNetConn 协议的抽象
	TcpNetConn struct {
		conn   net.Conn   // 连接
		config *NetConfig // 配置
	}

	// KcpNetConn 协议的抽象
	KcpNetConn struct {
		conn   net.Conn
		config *NetConfig // 配置
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
func NewTcpNetConn(c net.Conn, config *NetConfig) (NetConn, error) {
	return &TcpNetConn{conn: c, config: config}, nil
}

// Accept 接收连接
func (l *TcpNetListener) Accept() (NetConn, error) {
	c, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	if l.config.TLSConfig != nil {
		tlsConn := tls.Server(c, l.config.TLSConfig)

		return NewTcpNetConn(tlsConn, l.config)
	}

	return NewTcpNetConn(c, l.config)
}

// Close 关闭连接监听器
func (l *TcpNetListener) Close() {
	l.listener.Close()
}

// ReadPacket 读取数据报文
func (t *TcpNetConn) ReadPacket() (Packet, error) {
	t.conn.SetReadDeadline(time.Now().Add(t.config.ReadeTimeout))
	return Packing(t.conn)
}

// WritePacket 写入数据报文
func (t *TcpNetConn) WritePacket(p Packet) error {
	t.conn.SetWriteDeadline(time.Now().Add(t.config.WriteTimeout))
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
	kcpConn.SetNoDelay(1, 10, 2, 1)
	kcpConn.SetStreamMode(true)
	kcpConn.SetWindowSize(4096, 4096)
	kcpConn.SetReadBuffer(4 * 65536 * 1024)
	kcpConn.SetWriteBuffer(4 * 65536 * 1024)
	kcpConn.SetACKNoDelay(true)

	if l.config.TLSConfig != nil {
		tlsConn := tls.Server(kcpConn, l.config.TLSConfig)
		return NewKcpNetConn(tlsConn, l.config)
	}

	return NewKcpNetConn(kcpConn, l.config)
}

// Close 关闭连接监听器
func (l *KcpNetListener) Close() {
	l.listener.Close()
}

// NewKcpNetConn 创建一个kcp连接
func NewKcpNetConn(c net.Conn, config *NetConfig) (NetConn, error) {
	return &KcpNetConn{conn: c, config: config}, nil
}

// ReadPacket 读取数据包
func (k *KcpNetConn) ReadPacket() (Packet, error) {
	k.conn.SetReadDeadline(time.Now().Add(k.config.ReadeTimeout))
	return Packing(k.conn)
}

// WritePacket 写入数据包
func (k *KcpNetConn) WritePacket(p Packet) error {
	k.conn.SetWriteDeadline(time.Now().Add(k.config.WriteTimeout))
	_, err := k.conn.Write(p.Serialize())
	return err
}

// Close 关闭连接
func (k *KcpNetConn) Close() {
	k.conn.Close()
}
