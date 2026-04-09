package service

import (
	accountv1 "Book_Exp/webook/api/proto/gen/account/v1"
	paymentv1 "Book_Exp/webook/api/proto/gen/payment/v1"
	"Book_Exp/webook/pkg/logger"
	"Book_Exp/webook/reward/domain"
	"Book_Exp/webook/reward/repositroy"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type WechatNativeRewardService struct {
	client paymentv1.WechatPaymentServiceClient
	repo   repository.RewardRepository
	l      logger.LoggerV1
	acli   accountv1.AccountServiceClient
}

func (w *WechatNativeRewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	// 先查询缓存，确认是否已经创建过了打赏的预支付订单
	codeUrl, err := w.repo.GetCachedCodeURL(ctx, r)
	if err == nil {
		return codeUrl, nil
	}
	r.Status = domain.RewardStatusInit
	rid, err := w.repo.CreateReward(ctx, r)
	resp, err := w.client.NativePrePay(ctx, &paymentv1.PrePayRequest{
		Amt: &paymentv1.Amount{
			Total:    r.Amt,
			Currency: "CNY",
		},
		// 想办法拼接出来一个 biz_trade_id
		BizTradeNo:  w.bizTradeNO(rid),
		Description: fmt.Sprintf("打赏-%s", r.Target.BizName),
	})
	if err != nil {
		return domain.CodeURL{}, err
	}
	cu := domain.CodeURL{
		Rid: rid,
		URL: resp.CodeUrl,
	}
	err1 := w.repo.CachedCodeURL(ctx, cu, r)
	if err1 != nil {
		w.l.Error("缓存二维码失败", logger.Error(err1))
	}
	return cu, err

}

func (w *WechatNativeRewardService) GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error) {
	//快路径
	r, err := w.repo.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, err
	}
	if r.Uid != uid {
		return domain.Reward{}, errors.New("查询的打赏记录和打赏人对不上")
	}
	// 已经是完结状态
	if r.Completed() {
		return r, nil
	}
	// 这个时候，考虑到支付到查询结果，我们搞一个慢路径
	resp, err := w.client.GetPayment(ctx, &paymentv1.GetPaymentRequest{BizTradeNo: w.bizTradeNO(r.Id)})
	if err != nil {
		// 这边我们直接返回从数据库查询的数据
		w.l.Error("慢路径查询支付结果失败",
			logger.Int64("rid", r.Id), logger.Error(err))
		return r, nil
	}
	// 更新状态
	switch resp.Status {
	case paymentv1.PaymentStatus_PaymentStatusFailed:
		r.Status = domain.RewardStatusFailed
	case paymentv1.PaymentStatus_PaymentStatusInit:
		r.Status = domain.RewardStatusInit
	case paymentv1.PaymentStatus_PaymentStatusSuccess:
		r.Status = domain.RewardStatusPayed
	case paymentv1.PaymentStatus_PaymentStatusRefund:
		// 理论上来说不可能出现这个，直接设置为失败
		r.Status = domain.RewardStatusFailed
	}
	err = w.repo.UpdateStatus(ctx, rid, r.Status)
	if err != nil {
		w.l.Error("更新本地打赏状态失败",
			logger.Int64("rid", r.Id), logger.Error(err))
		return r, nil
	}
	return r, nil
}
func (s *WechatNativeRewardService) bizTradeNO(rid int64) string {
	return fmt.Sprintf("reward-%d", rid)
}

func (w *WechatNativeRewardService) UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {
	rid := w.toRid(bizTradeNO)
	err := w.repo.UpdateStatus(ctx, rid, status)
	if err != nil {
		return err
	}
	// 完成了支付，准备入账
	if status == domain.RewardStatusPayed {
		r, err := w.repo.GetReward(ctx, rid)
		if err != nil {
			return err
		}
		// webook 抽成
		weAmt := int64(float64(r.Amt) * 0.1)
		_, err = w.acli.Credit(ctx, &accountv1.CreditRequest{
			Biz:   "reward",
			BizId: rid,
			Items: []*accountv1.CreditItem{
				{
					AccountType: accountv1.AccountType_AccountTypeReward,
					// 虽然可能为 0，但是也要记录出来
					Amt:      weAmt,
					Currency: "CNY",
				},
				{
					Account:     r.Uid,
					Uid:         r.Uid,
					AccountType: accountv1.AccountType_AccountTypeReward,
					Amt:         r.Amt - weAmt,
					Currency:    "CNY",
				},
			},
		})
		if err != nil {
			w.l.Error("入账失败了，快来修数据啊！！！",
				logger.String("biz_trade_no", bizTradeNO),
				logger.Error(err))
			// 做好监控和告警，这里
			return err
		}
	}
	return nil
}
func (s *WechatNativeRewardService) toRid(tradeNO string) int64 {
	ridStr := strings.Split(tradeNO, "-")
	val, _ := strconv.ParseInt(ridStr[1], 10, 64)
	return val
}
func NewWechatNativeRewardService(
	client paymentv1.WechatPaymentServiceClient,
	repo repository.RewardRepository,
	l logger.LoggerV1,
	acli accountv1.AccountServiceClient,
) RewardService {
	return &WechatNativeRewardService{client: client, repo: repo, l: l, acli: acli}
}
