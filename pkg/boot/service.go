package boot

import (
	"fmt"
	"net"

	"github.com/trainking/goboot/pkg/etcdx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// BaseService service的基础服务
type BaseService struct {
	BaseInstance

	Name string // 服务名称

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

	// 验证选项
	switch s.Config.GetString(serviceConfigKey("Authentication.Mode", s.Name)) {
	case "TLS":
		creds, err := credentials.NewServerTLSFromFile(s.Config.GetString(serviceConfigKey("Authentication.CertFile", s.Name)), s.Config.GetString(serviceConfigKey("Authentication.KeyFile", s.Name)))
		if err != nil {
			return err
		}
		s.GrpcServer = grpc.NewServer(grpc.Creds(creds))
	default:
		s.GrpcServer = grpc.NewServer()
	}

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

// serviceConfigKey 获取指定服务的配置
func serviceConfigKey(k string, name string) string {
	return fmt.Sprintf("%s."+k, name)
}
