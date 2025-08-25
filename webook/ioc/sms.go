package ioc

import (
	"Book_Exp/webook/internal/service/sms"
	"Book_Exp/webook/internal/service/sms/memory"
)

func InitSmsService() sms.Service {
	//换内存还是换别的 实现
	return memory.NewService() //
}
