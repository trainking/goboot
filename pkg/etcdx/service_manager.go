package etcdx

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/trainking/goboot/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

var (
	// DefaultMetaData 默认元数据
	DefaultMetaData = &MetaData{
		Network: "http",
	}
)

// ServiceManager 服务节点管理器
type ServiceManager struct {
	em      endpoints.Manager
	xClient *ClientX

	leaseID  clientv3.LeaseID // 租约ID
	target   string           // 服务名
	leaseTTL int64            // 租约的过期时间
	heartT   int              // 心跳维持时间,单位秒
	ctx      context.Context
	cancel   context.CancelFunc
}

// MetaData 元数据结构
type MetaData struct {
	Network string `json:"Network"`
}

// NewServiceManager 增加服务节点管理器
func NewServiceManager(xClient *ClientX, target string, leaseTTl int64, heartT int) (*ServiceManager, error) {
	target = strings.TrimRight(target, "/")

	// 设置Context控制租约过期
	ctx, cancel := context.WithCancel(context.Background())

	em, err := endpoints.NewManager(xClient.client, target)
	if err != nil {
		cancel()
		return nil, err
	}
	return &ServiceManager{
		em:       em,
		xClient:  xClient,
		target:   target,
		leaseTTL: leaseTTl,
		heartT:   heartT,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Register 注册到节点
func (s *ServiceManager) Register(addr string, metadate ...interface{}) error {
	var _metadata interface{}
	if len(metadate) > 0 {
		_metadata = metadate[0]
	} else {
		_metadata = DefaultMetaData
	}

	lease := clientv3.NewLease(s.xClient.client)
	leaseResp, err := lease.Grant(s.ctx, s.leaseTTL)
	if err != nil {
		return err
	}

	// 维持租约
	go func() {
		for {
			if _, err := lease.KeepAliveOnce(s.ctx, leaseResp.ID); err != nil {
				fmt.Printf("lease keep %s %s Error: %v\n", s.target, addr, err)
				return
			}

			time.Sleep(time.Duration(s.heartT) * time.Second)
		}
	}()

	s.leaseID = leaseResp.ID

	return s.PushEndpoint(addr, _metadata)
}

// Watch 检查数据的变更
func (s *ServiceManager) Watch(h func(key string, ep endpoints.Endpoint)) error {
	wChan, err := s.em.NewWatchChannel(s.ctx)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			e := recover()
			if e != nil {
				log.Errorf("ServiceManager.Watch Error: %v", e)
			}
		}()
		for u := range wChan {
			for _, ud := range u {
				h(ud.Key, ud.Endpoint)
			}
		}
	}()

	return nil
}

// List 获取当前target的所有节点
func (s *ServiceManager) List() (map[string]endpoints.Endpoint, error) {
	return s.em.List(s.ctx)
}

// PushEndpoint push节点数据
func (s *ServiceManager) PushEndpoint(addr string, metadata interface{}) error {
	ep := endpoints.Endpoint{
		Addr:     addr,
		Metadata: metadata,
	}
	return s.em.AddEndpoint(s.xClient.client.Ctx(), s.target+"/"+addr, ep, clientv3.WithLease(s.leaseID))
}

// Destory 销毁
func (s *ServiceManager) Destory(addr string) error {
	if err := s.em.DeleteEndpoint(s.xClient.client.Ctx(), s.target+"/"+addr); err != nil {
		return err
	}

	if s.cancel != nil {
		s.cancel()
	}
	return nil
}
