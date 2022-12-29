package server

import (
	"context"

	"github.com/trainking/goboot/internal/pb"
)

// GetUserInfo 获取用户信息
func (s *Server) GetUserInfo(ctx context.Context, args *pb.GetUserInfoArgs) (*pb.GetUserInfoReply, error) {

	return &pb.GetUserInfoReply{
		UserName: "hw",
	}, nil
}
