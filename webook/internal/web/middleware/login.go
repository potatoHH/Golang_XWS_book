package middleware

import (
	"encoding/gob"
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
	gob.Register(time.Now()) // 注册类型以便在 gob 编码/解码过程中能够正确处理自定义类型
	return func(ctx *gin.Context) {
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
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
		sess.Set("userId", id)         // 刷新session
		sess.Options(sessions.Options{ // 刷新session的过期时间
			//Secure:   true,      // https  开发环境不要用
			//HttpOnly: true,      // js无法访问
			MaxAge: 60,
		})
		now := time.Now()      //毫秒
		if updatetime == nil { //更新时间
			//没有刷新过,刚登陆,还没有刷新
			sess.Set("update_time", now)

			sess.Save()
			return
		}
		//update_time是存在的
		updatetimeVal, ok := updatetime.(time.Time)
		if !ok {
			//没有登录
			ctx.AbortWithStatus(http.StatusInternalServerError) //500
			return

		}
		if now.Sub(updatetimeVal) > time.Minute*30 {
			sess.Set("update_time", now)
			sess.Save()
		}
	}

}
