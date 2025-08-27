package web

import (
	"Book_Exp/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
}

func NewOAuth2WechatHandler(svc wechat.Service) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc: svc,
	}

}

func (h *OAuth2WechatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthUrl)
	//这边any万无一失
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthUrl(ctx *gin.Context) {
	url, err := h.svc.AuthUrl(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URl失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {

}
