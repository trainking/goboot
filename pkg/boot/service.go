package boot

import (
	"fmt"
	"skyland-dev-peripheral/pkg/etcdx"
)

type BaseServcie struct {
	BaseInstance

	Prefix string // 服务前缀
	Name   string // 服务名称

	serviceManager *etcdx.ServiceManager
}

func (s *BaseServcie) Init() error {
	s.BaseInstance.Init()

	xClient, err := etcdx.New(s.Config.GetStringSlice(fmt.Sprintf("%s.Etcd", s.Name)))
	if err != nil {
		return err
	}

	s.serviceManager, err = etcdx.NewServiceManager(xClient, fmt.Sprintf("%s/%s", s.Prefix, s.Name), 15, 10)
	if err != nil {
		return err
	}

	return nil
}

func (s *BaseServcie) Start() error {
	return s.serviceManager.Register(s.Addr)
}

func (s *BaseServcie) Stop() {
	if s.serviceManager != nil {
		s.serviceManager.Destory(s.Addr)
	}
}
