package service

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// RequestIDKey requestID的键
const RequestIDKey = "X-Request-Id"

// WithRequestIDContext 加入requestID的context
func WithRequestIDContext(ctx context.Context, requestID string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, RequestIDKey, requestID)
}

// GetRequestID 获取RequestID
func GetRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		requestID := md.Get(RequestIDKey)
		if len(requestID) > 0 {
			return requestID[0]
		}
	}
	return ""
}
