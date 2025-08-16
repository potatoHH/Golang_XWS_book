package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signNmae *string
}

func NewService(client *sms.Client, appId, signName string) *Service {
	return &Service{
		client:   client,
		appId:    ekit.ToPtr[string](appId),    // 字符串转指针
		signNmae: ekit.ToPtr[string](signName), // 字符串转指针
	}
}

func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest() // 创建发送短信的请求
	req.SmsSdkAppId = s.appId      // 应用ID
	req.SignName = s.signNmae      // 签名
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) == "Ok" {
			return fmt.Errorf("短信发送失败 %s,原因是:%s", *status.Code, *status.Message)
		}
	}
	return nil

	//req.PhoneNumberSet = slice.Map[string, *string](numbers, func(idx int, src string) *string { // 手机号
	//	return &src
	//})
	//req.TemplateParamSet = slice.Map[string, *string](args, func(idx int, src string) *string { // 模板参数
	//	return &src
	//})
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
