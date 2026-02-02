package ioc

import (
	"Book_Exp/webook/internal/repository/dao"
	"Book_Exp/webook/pkg/logger"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13316)/webook",
	}
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		////缺了一个writer
		//Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
		//	SlowThreshold:             time.Microsecond * 10, //50/100毫米            // 慢查询阈值
		//	IgnoreRecordNotFoundError: true,                  // 忽略记录未找到错误
		//	ParameterizedQueries:      true,                  // 使用参数化查询
		//	LogLevel:                  glogger.Info,          // 日志级别
		//}),
	})
	if err != nil {
		panic(err)
	}
	//vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
	//	Namespace: "geekbing",
	//	Name:      "gorm_query_time",
	//	Subsystem: "webook",
	//	Help:      "统计GORM的执行时间",
	//	ConstLabels: map[string]string{
	//		"db": "webook",
	//	},
	//	Objectives: map[float64]float64{
	//		0.5:   0.01,
	//		0.75:  0.01,
	//		0.9:   0.01,
	//		0.95:  0.01,
	//		0.99:  0.01,
	//		0.999: 0.01,
	//	},
	//},
	//
	//	//如果join查询,talbe 就是 join在一起  ,或者table就是主表,A join B 记录的是A
	//	[]string{"type", "table"})
	//pcb := &Callbacks{
	//	vector: vector,
	//}
	//pcb.RegisterAll(db)
	////插件用法
	//db.Use(pcb)
	//prometheus.MustRegister(vector)
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type Callbacks struct {
	vector *prometheus.SummaryVec
}

func (pcb *Callbacks) Name() string {
	return "prometheus_query"
}

func (pcb *Callbacks) Initialize(db *gorm.DB) error {
	//TODO implement me
	panic("implement me")
}

func (pcb *Callbacks) RegisterAll(db *gorm.DB) {
	//监控查询的时间
	err := db.Callback().Create().Before("*").
		Register("prometheus_create_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", pcb.after("create"))
	err = dao.InitTable(db) // 初始化表
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", pcb.after("update"))
	err = dao.InitTable(db) // 初始化表
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_delete_after", pcb.after("delete"))
	err = dao.InitTable(db) // 初始化表
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().Before("*").
		Register("prometheus_raw_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", pcb.after("raw"))
	err = dao.InitTable(db) // 初始化表
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().Before("*").
		Register("prometheus_row_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").
		Register("prometheus_row_after", pcb.after("row"))
	err = dao.InitTable(db) // 初始化表
	if err != nil {
		panic(err)
	}
}

func (c *Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}
func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		c.vector.WithLabelValues(typ, table).Observe(float64(time.Since(startTime).Milliseconds()))
	}
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
