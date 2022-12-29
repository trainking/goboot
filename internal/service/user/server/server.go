package server

import (
	"net"

	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/log"
	"github.com/trainking/goboot/pkg/utils"
	"google.golang.org/grpc"
)

const ServicePrefix = "/services/"

type (
	Server struct {
		pb.UnimplementedUserServiceServer
		boot.BaseServcie
	}
)

// New 创建Server， 使用配置文件
func New(name string, configPath string, addr string, instanceId int64) boot.Instance {
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	server := new(Server)
	server.Prefix = ServicePrefix
	server.Name = name
	server.Config = v
	server.Addr = addr
	server.IntanceID = instanceId
	return server
}

// Init 服务初始化
func (s *Server) Init() error {
	var err error
	if err = s.BaseServcie.Init(); err != nil {
		log.Errorf("BaseServcie error %v", err)
		return err
	}

	return nil
}

// Start 开启服务
func (s *Server) Start() error {
	if err := s.BaseServcie.Start(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	grpcS := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcS, s)

	return grpcS.Serve(lis)
}
