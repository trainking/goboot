package gameapi

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Client 客户端实现，用于压测和机器人测试
type Client struct {
	Conn      net.Conn
	closeConn chan struct{}
	closeOnce sync.Once
	waitGroup sync.WaitGroup

	running bool

	receiveChan chan Packet
	sendChan    chan Packet
}

// NewClient 新客户端
// conn 连接协议实例
// readLimit 最大读取包
// sendLimit 最大写入包
// heart 心跳周期
func NewClient(conn net.Conn, readLimit int, sendLimit int, heart time.Duration) *Client {
	client := &Client{
		Conn:        conn,
		closeConn:   make(chan struct{}),
		receiveChan: make(chan Packet, readLimit),
		sendChan:    make(chan Packet, sendLimit),
		running:     true,
	}

	client.asyncDo(client.readLoop)
	client.asyncDo(client.sendLoop)

	client.KeepAlive(heart)

	return client
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
	if _, err := c.Conn.Write(HeartPacket.Serialize()); err != nil {
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

		n, err := Packing(c.Conn)
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
			if _, err := c.Conn.Write(p.Serialize()); err != nil {
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
