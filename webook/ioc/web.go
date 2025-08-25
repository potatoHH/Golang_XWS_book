package ioc

import (
	"Book_Exp/webook/internal/web"
	"Book_Exp/webook/internal/web/middleware"
	ratelimit "Book_Exp/webook/pkg/ginx/middlewares/ratlimit"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitGin(middl []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middl...)
	hdl.RegisterRoutes(server)
	return server
}
func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewLoginJwtMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("users/login_sms/code/send").
			IgnorePaths("users/login_sms").
			Build(),
		ratelimit.NewBuilder(redisClient, time.Minute, 100).Build(), //在多长时间内允许多少请求
	}
}

func corsHdl() gin.HandlerFunc {
	//跨域
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Content-Type", "authorization"},
		//暴露给前端的header
		ExposeHeaders: []string{"x-jwt-token"},
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
