package ciruitbreaker

import (
	"context"

	"github.com/go-kratos/aegis/circuitbreaker"

	"google.golang.org/grpc"
)

type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
	//设置标志位
}

// TODO 熔断器
func (s *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if s.breaker.Allow() == nil {
			resp, err = handler(ctx, req)
			if err != nil {
				//没有区别业务错误和系统错误
				s.breaker.MarkFailed()
			} else {
				s.breaker.MarkSuccess()
			}
		}
		s.breaker.MarkFailed()
		//触发了熔断器
		return nil, err

	}

}
