package ioc

import (
	"Book_Exp/webook/internal/service/sms"
	"Book_Exp/webook/internal/service/sms/memory"

	"github.com/redis/go-redis/v9"
)

func InitSmsService(cmd redis.Cmdable) sms.Service {
	//换内存还是换别的 实现
	return memory.NewService() //
}
