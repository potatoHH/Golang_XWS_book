package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	AtKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RfKey = []byte("95erj3fUD7fo0mlYdDbncXz4VD2igvh2")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}
func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(ctx, uid, ssid)
	return err
}

func (h *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	//我们使用jwt登录校验
	tokenHeader := ctx.GetHeader("Authorization") //获取token
	segs := strings.Split(tokenHeader, " ")       // Authorization: Bearer xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

	if len(segs) != 2 { // 判断 Bearer
		return ""
	}
	return segs[1]
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	//返回前端错误信息
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	claims := ctx.MustGet("user").(UserClaims)
	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}

func (h *RedisJWTHandler) CheackSession(ctx *gin.Context, ssid string) error {
	val, err := h.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if val == 0 {
			return nil
		}
		return errors.New("session 已经无效")
	default:
		return err
	}
}

func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		//设置过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //设置过期时间
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(), // 设置用户代理
	}
	//使用JWT设置登录状态  比如要求userid放入token中
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) // 创建一个token 使用jwt
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil

}
func (h *RedisJWTHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	rc := RefreshClamis{
		Uid:  uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			//设置过期时间7天
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, rc) // 创建一个token
	refreshTokenStr, err := refreshToken.SignedString(RfKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil

}
