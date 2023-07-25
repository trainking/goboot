package gameapi

import (
	"encoding/json"
	"fmt"

	"github.com/trainking/goboot/pkg/etcdx"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
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
		ID       int64       `json:"id"`       // 实例ID
		Version  string      `json:"version"`  // 当前服务的版本
		Network  string      `json:"network"`  // 传输协议
		UseTLS   bool        `json:"use_tls"`  // 是否启用TLS
		Fuse     bool        `json:"fuse"`     // 熔断开关，true为开启，熔断状态下，服务不接受新连接，等待服务器降到0
		State    int         `json:"state"`    // 状态
		Password string      `json:"password"` // 该服务加密使用的密码
		OutUrl   string      `json:"out_url"`  // 外部访问的URL
		Extra    interface{} `json:"extra"`    // 额外数据，自定义舒勇
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

	// 注册到服务中
	err = a.serviceManager.Register(a.Addr, a.gd)
	if err != nil {
		return err
	}
	// 监听值的变更
	err = a.serviceManager.Watch(func(key string, ep endpoints.Endpoint) {
		if ep.Addr == a.Addr {
			b, _ := json.Marshal(ep.Metadata)
			json.Unmarshal(b, &a.gd)
		}
	})
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

// UpdateGdExtra 修改服务的
func (a *App) UpdateGdExtra(extra interface{}) error {
	a.gd.Extra = extra

	return a.serviceManager.PushEndpoint(a.Addr, a.gd)
}

// EnableFuse 打开熔断开关
func (a *App) EnableFuse() error {
	a.gd.Fuse = true

	return a.serviceManager.PushEndpoint(a.Addr, a.gd)
}

// DisableFuse 关闭熔断开关
func (a *App) DisableFuse() error {
	a.gd.Fuse = false

	return a.serviceManager.PushEndpoint(a.Addr, a.gd)
}

// GetGd 获取此服务的元数据
func (a *App) GetGd() GameMetaData {
	return a.gd
}
