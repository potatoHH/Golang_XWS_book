package sms

import (
	"Book_Exp/webook/internal/repository"
	"Book_Exp/webook/internal/service/sms/tencent"
	"context"
	"fmt"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

const (
	codeTplId = "1877556"
)

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc tencent.Service
	//tplId string
}

func NewCodeService() *CodeService {
	return &CodeService{}
}

// 发送验证码
func (svc *CodeService) Send(ctx context.Context, phone string,
	//区别业务场景
	biz string) error {
	//3个步骤生成验证码    1生成验证码 2.塞进去redis ,3发送验证码
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	if err != nil {
		//这意味着redis有验证码, err 可能是超时了 err, 不知道发没发出去
	}
	return err

}

// 验证验证码  标准版
func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)

}

// 验证验证码  低频版
//func (svc *CodeService) VerifyV1(ctx context.Context, biz, phone, inputCode string) (bool, error) {
//
//}

// 生成验证码
func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000) //随机生成0-999999的数 包含
	//格式化, 不够的加上前导符
	return fmt.Sprintf("%6d", num)
}
