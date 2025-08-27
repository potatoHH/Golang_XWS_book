package wechat

import (
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"

	"net/url"
)

var redirectUrl = url.PathEscape("http://localhost:8080/oauth2/wechat/callback") //"http://localhost:8080/oauth2/wechat/callback" 需要编码
type Service interface {
	AuthUrl(ctx context.Context) (string, error)
}

type service struct {
	appId string
}

func NewService(appId string) Service {
	return &service{
		appId: appId,
	}

}

func (s *service) AuthUrl(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprint(urlPattern, s.appId, redirectUrl, state), nil
}
