package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	ExtractToken(ctx *gin.Context) string
	ClearToken(ctx *gin.Context) error
	CheackSession(ctx *gin.Context, ssid string) error
}

type UserClaims struct { // jwt
	jwt.RegisteredClaims // 继承
	//声明你自己要放进去token里面的数据
	Uid       int64
	UserAgent string
	Ssid      string
}
type RefreshClamis struct {
	jwt.RegisteredClaims
	Uid  int64
	Ssid string
}
