package gameapi

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/trainking/goboot/pkg/log"
)

const (
	// 向NATS中push推送给其他玩家消息的键
	NatsPushUserK = `GAME_PUSH_USER:%s`
)

type (
	// UserNats 向用户发送消息
	UserNats struct {
		serverName  string     // 服务器名
		nc          *nats.Conn // NATS消息分发
		receiveChan chan PushActorMessage
		sub         *nats.Subscription

		closeChan chan struct{}
		closeOnce sync.Once
	}

	// PushActorMessage 发送消息给远程玩家
	PushActorMessage struct {
		UserID int64  `json:"user_id"`
		OpCode uint16 `json:"opcode"`
		Msg    []byte `json:"msg"`
	}
)

// NewUserNats 创建一个向用户发送消息的nats
func NewUserNats(natsUrl string, serverName string) (*UserNats, error) {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}

	return &UserNats{nc: nc, serverName: serverName, closeChan: make(chan struct{})}, nil
}

// Puslish 发送用户消息
func (u *UserNats) Publish(userID int64, p Packet) error {
	pushMsg := PushActorMessage{
		UserID: userID,
		OpCode: p.OpCode(),
		Msg:    p.Body(),
	}

	b, err := json.Marshal(pushMsg)
	if err != nil {
		return err
	}

	return u.nc.Publish(fmt.Sprintf(NatsPushUserK, u.serverName), b)
}

// StartSubscribe 开始消费用户消息
func (u *UserNats) StartSubscribe(subLimit int) error {
	u.receiveChan = make(chan PushActorMessage)
	var subChan = make(chan *nats.Msg, subLimit)
	sub, err := u.nc.ChanSubscribe(fmt.Sprintf(NatsPushUserK, u.serverName), subChan)
	if err != nil {
		return err
	}
	u.sub = sub

	go func() {
		for {
			select {
			case <-u.closeChan:
				return
			case m := <-subChan:
				var msg PushActorMessage
				if err := json.Unmarshal(m.Data, &msg); err != nil {
					log.Errorf("PushActorMessage Unmarshal Error: %v", err)
					continue
				}

				u.receiveChan <- msg
			}
		}
	}()

	return nil
}

// Consume 读取消费到的用户消息
func (u *UserNats) Consume() <-chan PushActorMessage {
	if u.receiveChan == nil {
		panic("receive chan is nil")
	}

	return u.receiveChan
}

// Close 关闭
func (u *UserNats) Close() {
	u.closeOnce.Do(func() {
		close(u.closeChan)
		if u.receiveChan != nil {
			close(u.receiveChan)
		}
		if u.sub != nil {
			u.sub.Unsubscribe()
		}
		u.nc.Close()
	})
}
