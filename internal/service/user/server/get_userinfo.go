package server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/trainking/goboot/internal/pb"
	"github.com/trainking/goboot/pkg/log"
	"google.golang.org/grpc/metadata"
)

// GetUserInfo 获取用户信息
func (s *Server) GetUserInfo(ctx context.Context, args *pb.GetUserInfoArgs) (*pb.GetUserInfoReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		requestID := md.Get(echo.HeaderXRequestID)[0]
		log.Trace(requestID, "GetUserInfo", "rpc", "rpc incoming")
	}
	return &pb.GetUserInfoReply{
		UserName: "hw",
	}, nil
}
