package ratelimit

import (
	"Book_Exp/webook/internal/service/sms"
	"Book_Exp/webook/pkg/ratelimit"
	"context"
	"fmt"
)

type RateLimitSMSServiceV1 struct {
	sms.Service
	limiter ratelimit.Limiter //添加限流器
}

func NewRateLimitSMSServiceV1(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RateLimitSMSServiceV1{
		Service: svc,
		limiter: limiter,
	}
}

func (s *RateLimitSMSServiceV1) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "key")
	if err != nil {
		//TODO  系统错误,可以添加限流也可以不添加限流
		return fmt.Errorf("短信服务判断是否出现限流问题: %w", err)
	}
	if limited {
		return errLimited
	}
	//TODO 你在这里加 新特性的 代码
	err = s.Service.Send(ctx, tpl, args, numbers...) //短信发送
	//TODO 在这里也可以加 新特性
	return err
}
