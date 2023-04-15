package service

import (
	"github.com/spf13/viper"
	"github.com/trainking/goboot/pkg/etcdx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewGrpcClientConn 获取GrpcConn，根据配置
func NewGrpcClientConn(serviceName string, config *viper.Viper) (*grpc.ClientConn, error) {

	xClient, err := etcdx.New(config.GetStringSlice(serviceConfigKey("Etcd", serviceName)))
	if err != nil {
		return nil, err
	}

	var opts []grpc.DialOption
	switch config.GetString(serviceConfigKey("Authentication.Mode", serviceName)) {
	case "TLS":
		creds, err := credentials.NewClientTLSFromFile(config.GetString(serviceConfigKey("Authentication.CertFile", serviceName)), "")
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	default:
		opts = append(opts, grpc.WithInsecure())
	}

	return xClient.DialGrpc(config.GetString(serviceConfigKey("Prefix", serviceName))+"/"+serviceName, opts...)
}
