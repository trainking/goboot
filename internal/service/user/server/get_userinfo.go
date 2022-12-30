package server

import (
	"context"
	"fmt"

	"github.com/trainking/goboot/internal/pb"
	"google.golang.org/grpc/metadata"
)

// GetUserInfo 获取用户信息
func (s *Server) GetUserInfo(ctx context.Context, args *pb.GetUserInfoArgs) (*pb.GetUserInfoReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	fmt.Printf("MD: %v ok: %v \n", md, ok)
	return &pb.GetUserInfoReply{
		UserName: "hw",
	}, nil
}
