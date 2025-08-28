package web

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtHandler struct {
}

func (h jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {
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
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		return err
	}

	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

type UserClaims struct { // jwt
	jwt.RegisteredClaims // 继承
	//声明你自己要放进去token里面的数据
	Id        int64
	UserAgent string
}
