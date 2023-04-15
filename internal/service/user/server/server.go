package server

import (
	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/log"
	"github.com/trainking/goboot/pkg/service"
	"github.com/trainking/goboot/pkg/utils"
)

type (
	Server struct {
		pb.UnimplementedUserServiceServer
		service.BaseService
	}
)

// New 创建Server， 使用配置文件
func New(name string, configPath string, addr string, instanceId int64) boot.Instance {
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	server := new(Server)
	server.Name = name
	server.Config = v
	server.Addr = addr
	server.IntanceID = instanceId
	return server
}

// Init 服务初始化
func (s *Server) Init() error {
	var err error
	if err = s.BaseService.Init(); err != nil {
		log.Errorf("BaseServcie error %v", err)
		return err
	}

	pb.RegisterUserServiceServer(s.GrpcServer, s)

	return nil
}
