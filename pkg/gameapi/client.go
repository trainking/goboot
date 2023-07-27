package gameapi

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xtaci/kcp-go"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type (
	// Client tcp/kcp的传输客户端
	Client struct {
		Conn      NetConn
		closeConn chan struct{}
		closeOnce sync.Once
		waitGroup sync.WaitGroup

		running bool

		receiveChan chan Packet
		sendChan    chan Packet
	}
)

// NewClient 新客户端
// conn 连接协议实例
// readLimit 最大读取包
// sendLimit 最大写入包
// heart 心跳周期
func NewClient(network string, config NetConfig, readLimit int, sendLimit int, heart time.Duration) (*Client, error) {
	var netConn NetConn
	switch network {
	case "tcp":
		c, err := net.Dial("tcp", config.Addr)
		if err != nil {
			return nil, err
		}

		if config.TLSConfig != nil {
			tlsConn := tls.Client(c, config.TLSConfig)

			netConn = NewTcpNetConn(tlsConn, &config)
		} else {
			netConn = NewTcpNetConn(c, &config)
		}
	case "kcp":
		c, err := kcp.Dial(config.Addr)
		if err != nil {
			return nil, err
		}

		if config.TLSConfig != nil {
			tlsConn := tls.Client(c, config.TLSConfig)

			netConn = NewTcpNetConn(tlsConn, &config)
		} else {
			netConn = NewTcpNetConn(c, &config)
		}
	case "websocket":
		u := url.URL{Scheme: "ws", Host: config.Addr, Path: config.WebSocketPath}
		var dialer *websocket.Dialer
		if config.TLSConfig != nil {
			u.Scheme = "wss"
			dialer = &websocket.Dialer{TLSClientConfig: config.TLSConfig}
		} else {
			dialer = websocket.DefaultDialer
		}
		fmt.Println(u.String())
		c, _, err := dialer.Dial(u.String(), nil)
		if err != nil {
			return nil, err
		}
		netConn = NewWebSocketNetConn(c, &config)
	default:
		return nil, errors.New("no implement proto")
	}
	client := &Client{
		Conn:        netConn,
		closeConn:   make(chan struct{}),
		receiveChan: make(chan Packet, readLimit),
		sendChan:    make(chan Packet, sendLimit),
		running:     true,
	}

	client.asyncDo(client.readLoop)
	client.asyncDo(client.sendLoop)

	client.KeepAlive(heart)

	return client, nil
}

func (c *Client) asyncDo(fn func()) {
	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()
		fn()
	}()
}

// Ping 发送心跳
func (c *Client) Ping() {
	if err := c.Conn.WritePacket(HeartPacket); err != nil {
		c.Close()
		fmt.Printf("Ping error: %v\n", err)
	}
}

// KeepAlive 保持心跳
func (c *Client) KeepAlive(interval time.Duration) {
	if !c.running {
		return
	}
	if interval <= 0 {
		interval = time.Second * 3
	}

	time.AfterFunc(interval, func() {
		c.Ping()
		c.KeepAlive(interval)
	})
}

// readLoop 读取循环
func (c *Client) readLoop() {
	defer func() {
		recover()
		c.close()
	}()

	for {
		select {
		case <-c.closeConn:

			return
		default:
		}

		n, err := c.Conn.ReadPacket()
		if err != nil {
			return
		}

		if n.OpCode() > 0 {
			c.receiveChan <- n
		}
	}
}

// sendLoop 发送循环
func (c *Client) sendLoop() {
	defer func() {
		recover()
		c.close()
	}()

	for {
		select {
		case <-c.closeConn:
			return
		case p := <-c.sendChan:
			if err := c.Conn.WritePacket(p); err != nil {
				fmt.Println("Send Error: ", err)
				return
			}
		}
	}
}

// Send 发送消息
func (c *Client) Send(opcode uint16, msg protoreflect.ProtoMessage) error {
	p, err := CretaePbPacket(opcode, msg)
	if err != nil {
		return err
	}

	select {
	case <-c.closeConn:
		return errors.New("send close client")
	case c.sendChan <- p:
		return nil
	}
}

// Receive 读取消息
func (c *Client) Receive() <-chan Packet {
	return c.receiveChan
}

// close 关闭信号
func (c *Client) close() {
	c.closeOnce.Do(func() {
		close(c.closeConn)
		c.running = false
		c.Conn.Close()
	})
}

// Close 关闭服务
func (c *Client) Close() {
	c.close()

	c.waitGroup.Wait()
	close(c.receiveChan)
	close(c.sendChan)
}
