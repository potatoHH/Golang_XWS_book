package failover

import (
	"Book_Exp/webook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type FailoverSMSService struct {
	svcs []sms.Service //多个服务商
	idx  uint64
}

func NewFailoverSMSService(svcs ...sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil //发送成功了
		}
		//这边打印日志
		log.Println(err)
	}
	return errors.New("发送失败,所有的服务商都发送失败")
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//二话不说先把下标往后推一位
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ { // 如果i=9 length=5 那么9%5=4
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			//调用者设置的超时时间到了
			//调用者主动取消了
			return err
		}
		//其他情况打日志
		log.Println(err)
	}
	return errors.New("发送失败,所有服务商都发送失败")
}
