//go:build wireinject

package startup

import (
	"Book_Exp/webook/interactive/grpc"
	repository2 "Book_Exp/webook/interactive/repository"
	cache2 "Book_Exp/webook/interactive/repository/cache"
	dao2 "Book_Exp/webook/interactive/repository/dao"
	service2 "Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/repository/dao"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	dao2.NewGormInteractiveDAO,
	repository2.NewCachedInteractiveRepository,
	cache2.NewRedisInteractiveCache,
	dao.NewGormInteractiveDAO,
)

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service2.NewInteractiveService(nil, nil)
}
func InitInteractiverGRPCService() *grpc.InteractiveServiceServer {
	wire.Build(thirdProvider, interactiveSvcProvider, grpc.NewInteractiveServiceServer)
	return grpc.NewInteractiveServiceServer(nil)
}
