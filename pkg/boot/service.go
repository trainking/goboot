package boot

import (
	"fmt"
	"net"

	"github.com/trainking/goboot/pkg/etcdx"
	"google.golang.org/grpc"
)

// BaseService service的基础服务
type BaseService struct {
	BaseInstance

	Prefix string // 服务前缀
	Name   string // 服务名称

	Listener       net.Listener          // 网络句柄
	GrpcServer     *grpc.Server          // GRPC服务端
	serviceManager *etcdx.ServiceManager // etcd注册
}

func (s *BaseService) Init() error {
	var err error
	err = s.BaseInstance.Init()
	if err != nil {
		return err
	}

	xClient, err := etcdx.New(s.Config.GetStringSlice(fmt.Sprintf("%s.Etcd", s.Name)))
	if err != nil {
		return err
	}

	s.serviceManager, err = etcdx.NewServiceManager(xClient, fmt.Sprintf("%s/%s", s.Prefix, s.Name), 15, 10)
	if err != nil {
		return err
	}

	s.Listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.GrpcServer = grpc.NewServer()
	err = s.serviceManager.Register(s.Addr)

	return err
}

func (s *BaseService) Start() error {
	return s.GrpcServer.Serve(s.Listener)
}

func (s *BaseService) Stop() {
	if s.serviceManager != nil {
		s.serviceManager.Destory(s.Addr)
	}

	s.Listener.Close()
	s.GrpcServer.Stop()
}
