//go:build wireinject

package main

import (
	repository2 "Book_Exp/webook/interactive/repository"
	cache2 "Book_Exp/webook/interactive/repository/cache"
	service2 "Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/events/article"
	"Book_Exp/webook/internal/repository"
	article2 "Book_Exp/webook/internal/repository/article"
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao"
	article3 "Book_Exp/webook/internal/repository/dao/article"
	"Book_Exp/webook/internal/service"
	"Book_Exp/webook/internal/web"
	ijwt "Book_Exp/webook/internal/web/jwt"
	"Book_Exp/webook/ioc"

	"github.com/google/wire"
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitKafka,
		ioc.NewSyncProducer, ioc.NewConsumer,
		ioc.InitDB, ioc.InitRedis, //最基础的第三方依 赖
		ioc.InitLogger,
		dao.NewUserDao,
		ioc.InitGRPCClient,
		dao.NewGormInteractiveDAO,
		article3.NewArticleDao,

		cache.NewUserCache,
		cache.NewCodeCache,
		cache2.NewRedisInteractiveCache,
		cache.NewRedisArticleCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository2.NewCachedInteractiveRepository,

		article2.NewCacheArticleRepostiory,
		article.NewKafkaProducer,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		service2.NewInteractiveService,

		ioc.InitSmsService,
		ioc.InitOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		//直接基于内存的实现
		//memory.NewService,
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		//你的中间件呢,和注册路由
		//你这个地方没有用到前面的东西,
		ioc.InitGin,
		ioc.InitMiddlewares,
		ioc.NewWechatHandler,
		//consumer
		//article.NewInteractiveReadEventBatchConsumer,
		article.NewInteractiveReadEventBatchConsumer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
