//go:build wireinject

package main

import (
	"Book_Exp/webook/interactive/events"
	"Book_Exp/webook/interactive/grpc"
	"Book_Exp/webook/interactive/ioc"
	"Book_Exp/webook/interactive/repository"
	"Book_Exp/webook/interactive/repository/cache"
	"Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/repository/dao"

	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitKafka,
	ioc.InitGRPCxServer,
	ioc.NewConsumers,
	ioc.InitLogger,
)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGormInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitApp() *App {
	wire.Build(interactiveSvcProvider,
		thirdPartySet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
