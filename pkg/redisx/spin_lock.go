package redisx

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// DefaultSpinMutexConfig 默认一份自旋锁配置
var DefaultSpinMutexConfig = SpinMutexConfig{
	Retry:   3,
	Sleep:   1 * time.Second,
	Timeout: 5 * time.Second,
}

// SpinLock 自旋锁
type SpinMutex struct {
	conf SpinMutexConfig
	key  string // 锁的键

	redisClient redis.UniversalClient
}

// SpinMutexConfig 自旋锁配置
// * Retry 定义锁重试次数，到达次数还无法获得锁，则抛出retry to out
// * Sleep 每次重试后，阻塞时间，建议不要太长
// * Timout 持有锁的最长时间，即每次事务最长执行时间
type SpinMutexConfig struct {
	Retry   int           // 重试次数
	Sleep   time.Duration // 每次等待时间
	Timeout time.Duration // 持有锁最长时间
}

// NewSpinLock 创建一个重试锁
func NewSpinMutex(client redis.UniversalClient, key string, config SpinMutexConfig) *SpinMutex {
	return &SpinMutex{redisClient: client, key: key, conf: config}
}

// Lock 锁定
func (s *SpinMutex) Lock() error {
	for i := 0; i < s.conf.Retry; i++ {
		b, err := s.redisClient.SetNX(context.Background(), s.key, 1, s.conf.Timeout).Result()
		if err != nil {
			return err
		}
		// 未获取到锁
		if !b {
			time.Sleep(s.conf.Sleep)
			continue
		}
		return nil
	}
	return errors.New("retry to out")
}

// TryLock 尝试锁
func (s *SpinMutex) TryLock() bool {
	b, err := s.redisClient.SetNX(context.Background(), s.key, 1, s.conf.Timeout).Result()
	if err != nil {
		return false
	}
	return b
}

// Unlock 解锁
func (s *SpinMutex) Unlock() error {
	_, err := s.redisClient.Del(context.Background(), s.key).Result()
	return err
}
