//go:build wireinject

package main

import (
	"Book_Exp/webook/internal/repository"
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/service"
	"Book_Exp/webook/internal/web"
	"Book_Exp/webook/ioc"
	"Book_Exp/wire/repository/dao"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebService() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis, //最基础的第三方依 赖
		dao.NewUserDao,     //初始化dao
		cache.NewUserCache, //初始化缓存
		cache.NewCodeCache,
		repository.NewUserRepository,
		repository.NewCodeRepository,
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSmsService,
		//直接基于内存的实现
		//memory.NewService,
		web.NewUserHandler,
		//你的中间件呢,和注册路由
		//你这个地方没有用到前面的东西,
		ioc.InitWebService,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
