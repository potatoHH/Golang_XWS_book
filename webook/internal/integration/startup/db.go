package startup

import (
	"Book_Exp/webook/internal/repository/dao"
	"Book_Exp/webook/pkg/logger"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13316)/webook_default",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		//缺了一个writer
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			SlowThreshold:             time.Microsecond * 10, //50/100毫米            // 慢查询阈值
			IgnoreRecordNotFoundError: true,                  // 忽略记录未找到错误
			ParameterizedQueries:      true,                  // 使用参数化查询
			LogLevel:                  glogger.Info,          // 日志级别
		}),
	})

	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db) // 初始化表
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

//type Dosomting interface {
//	DoABC() string
//}
//type DosomtingFunc func() string //单方法
//
//func (d DosomtingFunc) DoABC() string {
//	return d()
//}
