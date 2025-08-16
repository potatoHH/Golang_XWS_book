package retryable

//import (
//	"Book_Exp/webook/internal/service/sms/tencent"
//	"context"
//)
//
//// 小心并发问题
//type Service struct {
//	svc tencent.Service
//	//重试次数
//	retryCnt int
//}
//
//func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
//	err := s.svc.Send(ctx, tplId, args, numbers...)
//	for err != nil && s.retryCnt < 10 {
//		s.svc.Send(ctx, tplId, args, numbers...)
//		s.retryCnt++
//	}
//	return err
//}
