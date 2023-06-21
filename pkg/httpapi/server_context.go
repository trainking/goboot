package httpapi

import (
	"errors"
	"sync"
)

type (
	// ServerContext 服务上下文，用于各个Module之间共享资源
	ServerContext struct {
		mm map[string]interface{}
		mu sync.RWMutex
	}
)

// newServerContext 创建一个ServerContext
func newServerContext() *ServerContext {
	return &ServerContext{
		mm: make(map[string]interface{}),
	}
}

// Add 加入一个共享的资源
func (s *ServerContext) Add(key string, val interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mm[key] = val
}

// Get 获取被共享的资源，如果key找不到，则返回Error
func (s *ServerContext) Get(key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RLock()
	if v, ok := s.mm[key]; ok {
		return v, nil
	}

	return nil, errors.New("key not found")
}
