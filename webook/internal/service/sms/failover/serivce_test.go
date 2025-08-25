package failover

import (
	"Book_Exp/webook/internal/service/sms"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) []sms.Service
		wantErr string
	}{
		{
			name: "发送成功",
			//mock: func(ctrl *gomock.Controller) []sms.Service {
			//	svc0 := failovermocks.NewMockService(ctrl)
			//	svc0.EXPECT().
			//		Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			//		Return(nil)
			//	return svc0
			//},
			//return []sms.Service{svc0}
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailoverSMSService(tc.mock(ctrl)...)
			err := svc.Send(context.Background(), "mytpl", []string{"123"}, "15432432432")

			assert.Equal(t, tc.wantErr, err)

		})
	}

}
