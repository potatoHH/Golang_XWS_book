package auth

import (
	"Book_Exp/webook/internal/service/sms"
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type SMSAuthService struct {
	svc sms.Service
	key string
}

func NewSMSAuthService(svcs sms.Service) sms.Service {
	return &SMSAuthService{
		svc: svcs,
	}

}

// TODO Send发送,其中biz必须是线下申请的一个代表业务方的token
func (s *SMSAuthService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	var tc Claims
	//权限校验
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token 不合法")
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type Claims struct {
	Tpl string
	jwt.RegisteredClaims
}
