package redisx

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"

	jsoniter "github.com/json-iterator/go"
)

type (
	// Cache 定义一个缓存的操作
	Cache interface {
		// Get 读取缓存， i 是读出的数据，必须是一个指针
		Get(ctx context.Context, i interface{}) error

		// Delete 删除缓存
		Delete(ctx context.Context) error

		// SetCodec 设置序列化协议
		SetCodec(c Codec)

		// SetLoadHandler 设置加载函数
		SetLoadHandler(h LoadHandler)
	}

	// Codec 是缓存结构序列化的抽象定义
	Codec interface {
		// Unmarshal 反序列化
		Unmarshal(data []byte, v interface{}) error

		// Marshal 序列化
		Marshal(v interface{}) ([]byte, error)
	}

	// LoadHandler 是 Get获取不到缓存时，应该调用的加载的函数
	LoadHandler func(context.Context) (interface{}, error)

	// StringCache 使用redis的string结构作缓存的结构体
	StringCache struct {
		client redis.UniversalClient // redis客户端
		ttl    time.Duration         // 缓存的过期时间
		key    string                // 缓存的键

		Codec       Codec       // 序列化协议
		LoadHandler LoadHandler // 加载函数

		sg singleflight.Group
	}
)

var (
	DefaultCodec = jsoniter.ConfigCompatibleWithStandardLibrary
)

// NewStringCache new一个StringCache，codec 默认使用
func NewStringCache(client redis.UniversalClient, key string, ttl time.Duration) Cache {
	sc := &StringCache{
		client: client,
		key:    key,
		ttl:    ttl,
	}

	// 使用默认的Codec
	sc.Codec = DefaultCodec

	return sc
}

// Get 读取缓存
func (s *StringCache) Get(ctx context.Context, i interface{}) error {
	_, err, _ := s.sg.Do(s.key, func() (interface{}, error) {

		id := reflect.ValueOf(i)
		if id.Kind() != reflect.Ptr {
			return nil, errors.New("i is not a pointer")
		}

		// fast，读取缓存
		v, err := s.client.Get(ctx, s.key).Result()
		if err == nil {
			err = s.Codec.Unmarshal([]byte(v), i)
			if err != nil {
				s.client.Del(ctx, s.key)
				return nil, err
			}
			return nil, nil
		}

		// slow， 加载缓存
		if err == redis.Nil {
			if s.LoadHandler == nil {
				return nil, errors.New("LoadHandler is nil")
			}

			d, err := s.LoadHandler(ctx)
			if err != nil {
				return nil, err
			}

			b, err := s.Codec.Marshal(d)
			if err != nil {
				return nil, err
			}

			_, err = s.client.Set(ctx, s.key, string(b), s.ttl).Result()
			if err != nil {
				return nil, err
			}

			s.Codec.Unmarshal(b, i)
			return nil, nil
		}
		return nil, err
	})
	return err
}

// Delete 删除缓存
func (s *StringCache) Delete(ctx context.Context) error {
	_, err, _ := s.sg.Do(s.key, func() (interface{}, error) {
		return s.client.Del(ctx, s.key).Result()
	})
	return err
}

// SetCodec 设置Codec
func (s *StringCache) SetCodec(c Codec) {
	s.Codec = c
}

// SetLoadHandler 设置LoadHandler
func (s *StringCache) SetLoadHandler(h LoadHandler) {
	s.LoadHandler = h
}
