package web

import (
	"Book_Exp/webook/internal/service"
	"Book_Exp/webook/internal/service/oauth2/wechat"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
)

type OAuth2WechatHandler struct {
	svc wechat.Service
	jwtHandler
	userSvc  service.UserServiceV1
	stateKey []byte
	cfg      WechatHandlerConfig
}
type WechatHandlerConfig struct {
	Secure bool
}

func NewOAuth2WechatHandler(svc wechat.Service, userSvc service.UserServiceV1, cfg WechatHandlerConfig) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:        svc,
		userSvc:    userSvc,
		stateKey:   []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvh1"),
		cfg:        cfg,
		jwtHandler: NewJWTHandler(),
	}

}

func (h *OAuth2WechatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthUrl)
	//这边any万无一失
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthUrl(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthUrl(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URl失败",
		})
		return
	}
	if err = h.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, StateCliams{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr, 600, "/oauth2/wechat/callback", "", h.cfg.Secure, true)
	return nil
}

type StateCliams struct {
	State string
	jwt.RegisteredClaims
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登陆失败",
		})
		return
	}
	info, err := h.svc.VerityCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//从userService里面拿id
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})

	}
	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	err = h.setRefreshToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "校验成功",
	})

}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	//校验state
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		//记录日志
		return fmt.Errorf("无法获取cookie,%w", err)
	}
	var sc StateCliams
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (any, error) {
		return h.stateKey, nil

	})
	if err != nil || !token.Valid {
		//记录日志
		return fmt.Errorf("cookie 不是合法的jwtCookie,%w", err)
	}
	//比较state
	if sc.State != state {
		//记录日志
		return errors.New("state 被篡改")

	}
	return nil
}
