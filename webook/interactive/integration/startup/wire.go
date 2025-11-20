package startup

import (
	"Book_Exp/webook/interactive/grpc"
	repository2 "Book_Exp/webook/interactive/repository"
	cache2 "Book_Exp/webook/interactive/repository/cache"
	service2 "Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/repository/dao"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	dao.NewGormInteractiveDAO,
	repository2.NewCachedInteractiveRepository,
	cache2.NewRedisInteractiveCache,
)

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service2.NewInteractiveService(nil, nil)
}
func InitInteractiverGRPCService() *grpc.InteractiveServiceServer {
	wire.Build(thirdProvider, interactiveSvcProvider, grpc.NewInteractiveServiceServer)
	return grpc.NewInteractiveServiceServer(nil)
}
