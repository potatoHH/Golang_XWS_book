package ratelimit

import (
	"testing"
)

func TestRateLimitSMSService_Send(t *testing.T) {
	//testCases := []struct {
	//	name string
	//	mock func(ctrl *gomock.Controller)
	//	ctx  context.Context
	//	ctrl *gomock.Controller
	//}{
	//	{
	//		name: "限流成功",
	//		mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
	//			svc :=
	//		},
	//	},
	//}
	//for _, tc := range testCases {
	//	t.Run(tc.name, func(t *testing.T) {
	//		ctrl := gomock.NewController(t)
	//		defer ctrl.Finish()
	//		l = NewRateLimitSMSService(tc.ctx, tc.mock(ctrl))
	//	})
	//}

}
