package ratelimit

import (
	grpc2 "Book_Exp/grpc"
	"Book_Exp/webook/pkg/logger"
	"Book_Exp/webook/pkg/ratelimit"
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InterceptorBuilder struct {
	limiter ratelimit.Limiter
	key     string
	l       logger.LoggerV1
}

// key user-service
// 整个集群限流   limiter:service:user:
// 用户服务限流   limiter:service:user:UserService
func NewInterceptorBuilder(limiter ratelimit.Limiter, key string, l logger.LoggerV1) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key, l: l}
}

func (b *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ok, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			b.l.Error("判定限流出现问题", logger.Error(err))
			//考虑你是用保守的还是激进的策略
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			//激进的策略
			//return handler(ctx, req)
		}
		//限流了才会返回true
		if ok {
			b.l.Error("触发限流", logger.String("key", b.key))
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		return handler(ctx, req)

	}

}

func (b *InterceptorBuilder) BuildClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ok, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			b.l.Error("判定限流出现问题", logger.Error(err))
			return status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		if ok {
			b.l.Error("触发限流", logger.Error(err))
		}
		return invoker(ctx, method, req, reply, cc, opts...)

	}

}

// 服务级别限流
func (b *InterceptorBuilder) BuilderServerInterceptorService() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		//prefix 这里做成参数
		if strings.HasPrefix(info.FullMethod, "/UserService") {
			ok, err := b.limiter.Limit(ctx, "limiter:service:user:UserService")
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

//灵活的业务限流

func (b *InterceptorBuilder) BuilderServerInterceptorBiz() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		if getById, ok := req.(*grpc2.GetByIdRequest); ok {
			key := fmt.Sprintf("limiter:user:get_by_id_:%d", getById.Id)
			ok, err := b.limiter.Limit(ctx, key)
			if err != nil {
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
			if ok {
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
		}
		return handler(ctx, req)

	}
}
