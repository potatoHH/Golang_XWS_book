package ioc

import (
	"Book_Exp/webook/internal/web"
	ijwt "Book_Exp/webook/internal/web/jwt"
	"Book_Exp/webook/internal/web/middleware"
	"Book_Exp/webook/pkg/ginx"
	logger2 "Book_Exp/webook/pkg/ginx/middlewares/logger"
	"Book_Exp/webook/pkg/logger"
	"context"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

func InitGin(middl []gin.HandlerFunc, userhdl *web.UserHandler, oauth2WechatHdl *web.OAuth2WechatHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middl...)
	userhdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	return server
}
func InitMiddlewares(redisClient redis.Cmdable, jwthdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		ginx.InitCounter(prometheus.CounterOpts{
			Namespace: "ginx",
			Subsystem: "webook",
			Name:      "ginx_http",
			Help:      "Http的业务错误码",
		}),
		logger2.NewBuilder(func(ctx context.Context, al *logger2.AccessLog) {
			l.Debug("HTTP请求", logger.Field{
				Key:   "al",
				Value: al,
			})
		}).AllowRespBody().AllowRepBody().Build(),
		middleware.NewLoginJwtMiddlewareBuilder(jwthdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/refresh_token").
			Build(),
		//ratelimit.NewBuilder(redisClient, 100).Build(), //在多长时间内允许多少请求
	}
}

func corsHdl() gin.HandlerFunc {
	//跨域
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Content-Type", "authorization"},
		//暴露给前端的header
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		//是否允许你带cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
