package client

import (
	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/internal/service/user/server"
	"github.com/trainking/goboot/pkg/etcdx"
	"google.golang.org/grpc"
)

func NewUserService(xClient *etcdx.ClientX) (pb.UserServiceClient, error) {
	conn, err := xClient.DialGrpc(server.ServicePrefix+"UserService", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return pb.NewUserServiceClient(conn), nil
}

func NewUserServiceByMap(setting map[string]interface{}) (pb.UserServiceClient, error) {
	var etcdGateway []string

	for _, gateway := range setting["etcd"].([]interface{}) {
		etcdGateway = append(etcdGateway, gateway.(string))
	}

	xClient, err := etcdx.New(etcdGateway)
	if err != nil {
		return nil, err
	}
	return NewUserService(xClient)
}
