package grpc

import (
	"Book_Exp/webook/pkg/logger"
	"Book_Exp/webook/pkg/ratelimit"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	limiter ratelimit.Limiter
	key     string
	l       logger.LoggerV1
}

func NewInterceptorBuilder(limiter ratelimit.Limiter, key string, l logger.LoggerV1) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key, l: l}
}

// 服务级别限流
func (b *InterceptorBuilder) BuilderServerInterceptorServiceBiz() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		//prefix 这里做成参数
		if idReq, ok := req.(*GetByIdRequest); ok {
			ok, err := b.limiter.Limit(ctx, fmt.Sprintf("limiter:service:%s:%d", info.FullMethod, idReq.Id))
			if err != nil {
				return nil, status.Errorf(codes.ResourceExhausted, "服务级别触发限流")
			}
			if ok {
				return nil, status.Errorf(codes.ResourceExhausted, "服务级别触发限流")
			}
		}
		return handler(ctx, req)

	}
}
