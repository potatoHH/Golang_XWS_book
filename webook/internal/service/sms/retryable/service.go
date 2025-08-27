package retryable

import (
	"Book_Exp/webook/internal/service/sms"
	"context"
	"errors"
)

// 小心并发问题
type Service struct {
	svc sms.Service
	//重试次数
	retryMax int
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	cnt := 1
	for err != nil && cnt < s.retryMax {
		s.svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New("重试发送短信失败")
}
