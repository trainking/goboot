package etcdx

import (
	"context"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

// ServiceManager 服务节点管理器
type ServiceManager struct {
	em      endpoints.Manager
	xClient *ClientX

	target   string // 服务名
	leaseTTL int64  // 租约的过期时间
	heartT   int    // 心跳维持时间,单位秒
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewServiceManager 增加服务节点管理器
func NewServiceManager(xClient *ClientX, target string, leaseTTl int64, heartT int) (*ServiceManager, error) {
	target = strings.TrimRight(target, "/")

	em, err := endpoints.NewManager(xClient.client, target)
	if err != nil {
		return nil, err
	}
	return &ServiceManager{
		em:       em,
		xClient:  xClient,
		target:   target,
		leaseTTL: leaseTTl,
		heartT:   heartT,
	}, nil
}

// Register 注册到节点
func (s *ServiceManager) Register(addr string, metadata ...interface{}) error {
	// 设置Context控制租约过期
	s.ctx, s.cancel = context.WithCancel(context.Background())

	ep := endpoints.Endpoint{
		Addr: addr,
	}

	if len(metadata) > 0 {
		ep.Metadata = metadata[0]
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

	return s.em.AddEndpoint(s.xClient.client.Ctx(), s.target+"/"+addr, ep, clientv3.WithLease(leaseResp.ID))
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
