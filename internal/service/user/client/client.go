package client

import (
	"github.com/spf13/viper"
	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/boot"
)

// NewUserService 创建UserServcie客户端
func NewUserService(serviceName string, config *viper.Viper) (pb.UserServiceClient, error) {
	conn, err := boot.NewGrpcClientConn(serviceName, config)
	if err != nil {
		return nil, err
	}
	return pb.NewUserServiceClient(conn), nil
}
