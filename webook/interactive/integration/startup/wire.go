package startup

import (
	repository2 "Book_Exp/webook/interactive/repository"
	cache2 "Book_Exp/webook/interactive/repository/cache"
	service2 "Book_Exp/webook/interactive/service"
	dao2 "Book_Exp/webook/internal/repository/dao"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog)
var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewCachedInteractiveRepository,
	dao2.NewGormInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service2.NewInteractiveService(nil, nil)
}
