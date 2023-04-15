package service

import (
	"fmt"
	"net"

	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/etcdx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// BaseService service的基础服务
type BaseService struct {
	boot.BaseInstance

	Listener       net.Listener          // 网络监听
	GrpcServer     *grpc.Server          // GRPC服务端
	serviceManager *etcdx.ServiceManager // etcd注册
}

func (s *BaseService) Init() error {
	var err error
	err = s.BaseInstance.Init()
	if err != nil {
		return err
	}

	xClient, err := etcdx.New(s.Config.GetStringSlice(serviceConfigKey("Etcd", s.Name)))
	if err != nil {
		return err
	}

	s.serviceManager, err = etcdx.NewServiceManager(xClient, fmt.Sprintf("%s/%s", s.Config.GetString(serviceConfigKey("Prefix", s.Name)), s.Name), 15, 10)
	if err != nil {
		return err
	}

	s.Listener, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	// 注册的元数据
	var metadate *etcdx.MetaData

	// 验证选项
	switch s.Config.GetString(serviceConfigKey("Authentication.Mode", s.Name)) {
	case "TLS":
		creds, err := credentials.NewServerTLSFromFile(s.Config.GetString(serviceConfigKey("Authentication.CertFile", s.Name)), s.Config.GetString(serviceConfigKey("Authentication.KeyFile", s.Name)))
		if err != nil {
			return err
		}
		s.GrpcServer = grpc.NewServer(grpc.Creds(creds))
		metadate = &etcdx.MetaData{
			Network: "https",
		}
	default:
		s.GrpcServer = grpc.NewServer()
		metadate = etcdx.DefaultMetaData
	}

	err = s.serviceManager.Register(s.Addr, metadate)

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

// serviceConfigKey 获取指定服务的配置
func serviceConfigKey(k string, name string) string {
	return fmt.Sprintf("%s."+k, name)
}
