package etcdx

import (
	"context"
	"errors"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ErrorNotFound 表示该key找不多对应的值，即不存在
var ErrorNotFound = errors.New("not found")

const (
	// DefaultDialTimout 默认连接超时时间
	DefaultDialTimout = 10 * time.Second
)

// ChangeHandler 发生变化的函数
type ChangeHandler func(key string, value []byte)

// ClientX 是封装的Etcd客户端操作
type ClientX struct {
	client *clientv3.Client
	kv     clientv3.KV
}

// New 新建一个ClientX实例
func New(etcdGateway []string) (*ClientX, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdGateway,
		DialTimeout: DefaultDialTimout,
	})

	if err != nil {
		return nil, err
	}

	return &ClientX{client: client, kv: clientv3.NewKV(client)}, nil
}

// Put 设置一个键值对
func (c *ClientX) Put(ctx context.Context, key, value string) error {
	_, err := c.kv.Put(ctx, key, value)
	return err
}

// Get 获取一个键值对
func (c *ClientX) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.kv.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	for _, kv := range resp.Kvs {
		return kv.Value, nil
	}
	return nil, ErrorNotFound
}

// Delete 删除一个指定的键
func (c *ClientX) Delete(ctx context.Context, key string) error {
	_, err := c.kv.Delete(ctx, key)
	return err
}

// GetWithPrefix 前缀获取一组键值，找不到 ErrorNotFound
func (c *ClientX) GetWithPrefix(ctx context.Context, key string) ([][]byte, error) {
	resp, err := c.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	length := len(resp.Kvs)
	if length == 0 {
		return nil, ErrorNotFound
	}

	rList := make([][]byte, length)
	j := 0
	for _, kv := range resp.Kvs {
		rList[j] = kv.Value
		j++
	}
	return rList, nil
}

// WatchWithPrefix 监听某一组简直对的变化
func (c *ClientX) WatchWithPrefix(ctx context.Context, key string, putHandler ChangeHandler, delHandler ChangeHandler) {
	wc := c.client.Watch(ctx, key, clientv3.WithPrefix())
	for v := range wc {
		for _, e := range v.Events {
			switch e.Type {
			case mvccpb.DELETE: // 删除触发
				delHandler(string(e.Kv.Key), e.Kv.Value)
			case mvccpb.PUT: // put触发
				putHandler(string(e.Kv.Key), e.Kv.Value)
			}
		}
	}
}
