package gameapi

import (
	"fmt"

	"github.com/trainking/goboot/pkg/etcdx"
)

const (
	StateZero   = iota // 服务初始状态，服务连接数为0
	StateIdle          // 服务空闲，优先获得连接
	StateAliave        // 服务正常运行，均衡获得连接
	StateBusy          // 服务繁忙，滞后获得连接
	StateFull          // 服务满载，不获得新连接
)

type (
	GameMetaData struct {
		ID      int64  `json:"id"`      // 实例ID
		Network string `json:"network"` // 传输协议
		UseTLS  bool   `json:"use_tls"` // 是否启用TLS
		Fuse    bool   `json:"fuse"`    // 熔断开关，true为开启，熔断状态下，服务不接受新连接，等待服务器降到0
		State   int    `json:"state"`   // 状态
	}
)

// registerEtcd 注册到Etcd中
func (a *App) registerEtcd() error {
	etcdArr := a.Config.GetStringSlice("Etcd")
	if len(etcdArr) == 0 {
		return nil
	}
	xClient, err := etcdx.New(etcdArr)
	if err != nil {
		return err
	}

	a.serviceManager, err = etcdx.NewServiceManager(xClient, fmt.Sprintf("%s/%s", a.Config.GetString("Prefix"), a.Name), 15, 10)
	if err != nil {
		return err
	}

	err = a.serviceManager.Register(a.Addr, a.gd)
	return err
}

// UpdateGdState 修改服务状态
func (a *App) UpdateGdState(state int) error {
	if a.gd.State == state {
		return nil
	}

	a.gd.State = state
	return a.serviceManager.PushEndpoint(a.Addr, a.gd)
}
