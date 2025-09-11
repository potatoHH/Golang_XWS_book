package wechat

import (
	"Book_Exp/webook/internal/domain"
	"Book_Exp/webook/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"net/url"
)

var redirectUrl = url.PathEscape("https://meoying.com/oauth2/wechat/callback") //"http://localhost:8080/oauth2/wechat/callback" 需要编码
type Service interface {
	AuthUrl(ctx context.Context, state string) (string, error)
	VerityCode(ctx context.Context, code string) (domain.WecahteInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
	logger    logger.LoggerV1
}

func NewService(appId string, appSecret string, l logger.LoggerV1) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		//依赖注入,但是没有完全注入
		client: http.DefaultClient,
		logger: l,
	}
}

func (s *service) AuthUrl(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprint(urlPattern, s.appId, redirectUrl, state), nil
}

func (s *service) VerityCode(ctx context.Context, code string) (domain.WecahteInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WecahteInfo{}, err
	}
	resp, err := s.client.Do(req) // 发送请求
	if err != nil {
		return domain.WecahteInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body) // 创建一个解码器
	var res Result
	err = decoder.Decode(&res) // 解析响应
	if err != nil {
		return domain.WecahteInfo{}, err
	}

	if res.Errcode != 0 {
		return domain.WecahteInfo{}, fmt.Errorf("微信返回错误响应,错误码:%d,错误信息:%s", res.Errcode, res.Errmsg)
	}
	return domain.WecahteInfo{
		OpenId:  res.Openid,
		UnionId: res.Unionid,
	}, nil
}

type Result struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	ExpiresIn    int64  `json:"expires_in"`    // 有效期
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	Openid       string `json:"openid"`        // 用户唯一标识
	Scope        string `json:"scope"`         // 授权范围
	Unionid      string `json:"unionid"`       // 用户在开放平台的唯一标识符
	Errcode      int64  `json:"errcode"`       // 错误码
	Errmsg       string `json:"errmsg"`        // 错误信息
}
