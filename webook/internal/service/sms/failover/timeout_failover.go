package failover

import (
	"Book_Exp/webook/internal/service/sms"
	"context"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs      []sms.Service
	idx       int32 //索引
	cnt       int32 //连续超时的个数
	threshold int32 //阈值,连续超时的个数超过这个数字,就要切换
}

func NewTimeoutFailoverSMSService(cnt, idx int32, svcs ...sms.Service) sms.Service {
	return &TimeoutFailoverSMSService{
		svcs: svcs,
		cnt:  cnt,
		idx:  idx,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		//触发切换,重新计算新的下标
		newIdx := (idx + 1) % int32(len(t.svcs))
		//CAS 操作失败,说明有人切换了,所以你这里不需要检测返回值
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//说明切换了
			atomic.StoreInt32(&t.cnt, 0)
		}
		//else 就是出现并发,别人换成功了
		//idx=newIdx
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	//当使用的是svc
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case nil:
		//没有任何错误,重置计数器
		atomic.StoreInt32(&t.cnt, 0)
	case context.DeadlineExceeded: //超时
		atomic.AddInt32(&t.cnt, 1)
	default:
		//如果是别的异常活动,我们保持不动
		return err
	}
	return err
}
