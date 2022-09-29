package etcdx

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// LeaseX 租约的拓展
type LeaseX struct {
	target   string // 租约的前缀
	value    string // 推送的值
	leaseTTL int64  // 租约的过期时间
	heartT   int    // 心跳维持时间,单位秒

	closeCh chan struct{}

	xClient *ClientX
}

// Update 租约自动续约携程
func (l *LeaseX) Update() {

	lease := clientv3.NewLease(l.xClient.client)
	var leaseID clientv3.LeaseID = 0

	go func() {
		for {
			select {
			case <-l.closeCh:
				// 关闭退出
				return
			default:
				if leaseID == 0 {
					leaseResp, err := lease.Grant(context.TODO(), l.leaseTTL)
					if err != nil {
						panic(err)
					}

					key := fmt.Sprintf("%s%d", l.target, leaseResp.ID)
					if _, err = l.xClient.kv.Put(context.TODO(), key, l.value, clientv3.WithLease(leaseResp.ID)); err != nil {
						panic(err)
					}
					leaseID = leaseResp.ID
				} else {
					if _, err := lease.KeepAliveOnce(context.TODO(), leaseID); err == rpctypes.ErrLeaseNotFound {
						leaseID = 0
						continue
					}
				}

				// 维持心跳续约
				time.Sleep(time.Duration(l.heartT) * time.Second)
			}
		}
	}()
}

// Close 关闭租约
func (l *LeaseX) Close() {
	close(l.closeCh)
}
