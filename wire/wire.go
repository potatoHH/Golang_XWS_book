//go:build wireinject

// 让wire注入这里的代码
package wire

import (
	"Book_Exp/wire/repository"
	"Book_Exp/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.Repository {
	//这个方法里面传入各个组件的初始化方法
	wire.Build(repository.NewRepository, dao.NewUserDao, InitDB())
	return new(repository.Repository)
}
