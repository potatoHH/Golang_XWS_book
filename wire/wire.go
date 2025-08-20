//go:build wireinject

// 让wire注入这里的代码
package wire

import (
	"Book_Exp/wire/repository"
	"Book_Exp/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	//TODO 我只在这里面声明我要用的各种东西,但是具体怎么构造,怎么编排排序
	//这个方法里面传入各个组件的初始化方法
	wire.Build(repository.NewUserRepository, dao.NewUserDao, InitDB)
	return new(repository.UserRepository)
}
