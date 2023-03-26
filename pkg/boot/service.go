package boot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/trainking/goboot/pkg/etcdx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	// 根据配置启用验证
	switch s.Config.GetString(fmt.Sprintf("%s.Authentication.Mode", s.Name)) {
	case "TLS":
		// 加载TLS证书
		cert, err := ioutil.ReadFile(s.Config.GetString(fmt.Sprintf("%s.Authentication.ServerCrt", s.Name)))
		if err != nil {
			return err
		}
		// 加载TLS密钥
		key, err := ioutil.ReadFile(s.Config.GetString(fmt.Sprintf("%s.Authentication.ServerKey", s.Name)))
		if err != nil {
			return err
		}

		// 加载CA证书
		ca, err := ioutil.ReadFile(s.Config.GetString(fmt.Sprintf("%s.Authentication.CaCrt", s.Name)))
		if err != nil {
			return err
		}

		// 创建证书池
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(ca)

		// 创建TLS配置
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{
				{
					Certificate: [][]byte{cert},
					PrivateKey:  key,
				},
			},
			ClientCAs:  certPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}

		s.GrpcServer = grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
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
