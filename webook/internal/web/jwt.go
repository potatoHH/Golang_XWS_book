package web

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtHandler struct {
	atKey []byte //accesstoken
	rfKey []byte //refreshTokenKey
}

func NewJWTHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		rfKey: []byte("95erj3fUD7fo0mlYdDbncXz4VD2igvh2"),
	}
}

func (h *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		//设置过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //设置过期时间
		},
		Id:        uid,
		UserAgent: ctx.Request.UserAgent(), // 设置用户代理
	}
	//使用JWT设置登录状态  比如要求userid放入token中
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // 创建一个token 使用jwt
	tokenStr, err := token.SignedString(h.atKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
func (h *jwtHandler) setRefreshToken(ctx *gin.Context, uid int64) error {
	rc := RefreshClamis{
		uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			//设置过期时间7天
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, rc) // 创建一个token
	refreshTokenStr, err := refreshToken.SignedString(h.rfKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil

}

func ExtractToken(ctx *gin.Context) string {
	//我们使用jwt登录校验
	tokenHeader := ctx.GetHeader("Authorization") //获取token
	segs := strings.Split(tokenHeader, " ")       // Authorization: Bearer xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	if len(segs) != 2 { // 判断 Bearer
		return ""
	}
	return segs[1]
}

type UserClaims struct { // jwt
	jwt.RegisteredClaims // 继承
	//声明你自己要放进去token里面的数据
	Id        int64
	UserAgent string
}
type RefreshClamis struct {
	jwt.RegisteredClaims
	uid int64
}
