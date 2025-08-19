package ioc

import (
	"Book_Exp/webook/config"
	"Book_Exp/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := config.Config.DB.DSN // 数据库连接
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db) // 初始化表
	if err != nil {
		panic(err)
	}
	return db
}
