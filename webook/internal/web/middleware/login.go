package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct { // 登录中间件
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := sessions.Default(ctx)
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
		if sess.Get("userId") == nil {
			//没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized) //401
			return
		}

		//刷新登录时间
		updatetime := sess.Get("update_time")
		now := time.Now().UnixMilli() //毫秒
		if updatetime == nil {        //更新时间
			//没有刷新过,刚登陆,还没有刷新
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		//update_time是存在的
		if updatetimeVal, ok := updatetime.(int64); ok {
			if now - updatetimeVal  > 30* 60 * 1000 { //30分钟
				sess.Set("update_time", now)
				sess.Save()
				return
 			}

			}
		}
	}

}
