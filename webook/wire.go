//go:build wireinject

package wire

import (
	"Book_Exp/webook/internal/repository"
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao"
	"Book_Exp/webook/internal/service"
	"Book_Exp/webook/internal/web"
	"Book_Exp/webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis, //最基础的第三方依 赖
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSmsService,
		ioc.InitOAuth2WechatHandler,
		//直接基于内存的实现
		//memory.NewService,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		//你的中间件呢,和注册路由
		//你这个地方没有用到前面的东西,
		ioc.InitGin,
		ioc.InitMiddlewares,
		ioc.NewWechatHandler,
	)
	return new(gin.Engine)
}
