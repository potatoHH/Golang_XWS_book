package wechat

import (
	"Book_Exp/webook/payment/domain"
	"Book_Exp/webook/payment/events"
	"Book_Exp/webook/payment/repository"
	"Book_Exp/webook/pkg/logger"
	"context"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/partnerpayments/native"
)

type NativePaymentService struct {
	svc                  *native.NativeApiService
	appID                string
	mchID                string
	notifyURL            string
	repo                 repository.PaymentRepository
	l                    logger.LoggerV1
	producer             events.Producer
	nativeCBTypeToStatus map[string]domain.PaymentStatus
}

func NewNativePaymentService(svc *native.NativeApiService, l logger.LoggerV1, repo repository.PaymentRepository, appid, mchid string) PaymentService {
	return &NativePaymentService{
		svc:   svc,
		l:     l,
		repo:  repo,
		appID: appid,
		mchID: mchid,
		// 一般来说，这个都是固定的，基本不会变的
		notifyURL: "http://wechat.meoying.com/pay/callback",
		nativeCBTypeToStatus: map[string]domain.PaymentStatus{
			"SUCCESS":  domain.PaymentStatusSuccess,
			"PAYERROR": domain.PaymentStatusFailed,
			"NOTPAY":   domain.PaymentStatusInit,
			"CLOSED":   domain.PaymentStatusFailed,
			"REVOKED":  domain.PaymentStatusFailed,
			"REFUND":   domain.PaymentStatusRefund,
			// 其它状态你都可以加
		},
	}
}

func (n *NativePaymentService) Prepay(ctx context.Context, pmt domain.Payment) (string, error) {
	//唯一索引冲突，业务方唤起了支付，但是没付，下一次再过来，应该换成 BizTradeNo
	err := n.repo.AddPayment(ctx, pmt)
	if err != nil {
		return "", err
	}
	resp, result, err := n.svc.Prepay(ctx, native.PrepayRequest{
		SpAppid:     core.String(n.appID),
		SpMchid:     core.String(n.mchID),
		Description: core.String(pmt.Description),
		//设置30分钟有效时间
		TimeExpire: core.Time(time.Now().Add(time.Second * 30)),
		//这个是有讲究的
		//选择1:业务方直接给欸我,我透传,我啥也不太干
		//选择2:业务方给他的业务表示,我自己生成一个
		//不管1,2 业务方都必须传给我们一个唯一标识
		OutTradeNo: core.String(pmt.BizTradeNO),
		NotifyUrl:  core.String(n.notifyURL),
		Amount: &native.Amount{
			Currency: core.String(pmt.Amt.Currency),
			Total:    core.Int64(pmt.Amt.Total),
		},
	})
	n.l.Debug("微信pay响应",
		logger.Field{"resp", resp},
		logger.Field{"result", result},
	)
	if err != nil {
		return "", err
	}
	return *resp.CodeUrl, err
}
