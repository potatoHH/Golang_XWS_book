package middleware

import (
	"Book_Exp/webook/internal/web"
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// jwt登录校验
type LoginJwtMiddlewareBuilder struct { // 登录中间件
	paths []string
}

func NewLoginJwtMiddlewareBuilder() *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{}
}
func (l *LoginJwtMiddlewareBuilder) IgnorePaths(paths string) *LoginJwtMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}
func (l *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now()) // 注册类型以便在 gob 编码/解码过程中能够正确处理自定义类型
	return func(ctx *gin.Context) {
		//前面的session设置过了，注释这里不可能为nil
		//if sess == nil {
		//	//没有登录
		//	ctx.AbortWithStatus(http.StatusUnauthorized) //401
		//	return
		//}
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" { //登录和注册接口，不需要登录校验
			return
		}

		//我们使用jwt登录校验
		tokenHeader := ctx.GetHeader("Authorization") //获取token
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized) //401 没登陆
		}
		segs := strings.SplitN(tokenHeader, " ", 2)                  // Authorization: Bearer xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		if len(segs) != 2 || !strings.EqualFold(segs[0], "Bearer") { // 判断 Bearer
			ctx.AbortWithStatus(http.StatusUnauthorized) //401
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		//parsewithClaims里面一定要传指针 拿到token
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) { // 解析token
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized) //401
			return
		}
		//err为nil, token不为nil
		if token == nil || !token.Valid || claims.Uid <= 0 { // token校验失败
			ctx.AbortWithStatus(http.StatusUnauthorized) //401
			return
		}
		//每10秒刷新一次
		now := time.Now()
		if claims.ExpiresAt.Time.Sub(now) < time.Second*50 {
			//刷新过期时间
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24)) //设置过期时间
			tokenStr, err = token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}
		ctx.Header("x-jwt-token", tokenStr)
		ctx.Set("claims", claims)
	}

}
