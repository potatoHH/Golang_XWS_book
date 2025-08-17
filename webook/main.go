package main

import (
	"Book_Exp/webook/config"
	"Book_Exp/webook/internal/repository"
	"Book_Exp/webook/internal/repository/cache"
	"Book_Exp/webook/internal/repository/dao"
	"Book_Exp/webook/internal/service"
	"Book_Exp/webook/internal/service/sms/memory"
	"Book_Exp/webook/internal/web"
	"Book_Exp/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	//"github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initWebServer() *gin.Engine { // 初始化web服务
	server := gin.Default() // 创建一个gin服务
	/// 跨域
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"Content-Type", "authorization"},
		//暴露给前端的header
		ExposeHeaders: []string{"x-jwt-token"},
		//是否允许你带cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//开发环境
				return true
			}
			return strings.Contains(origin, "yourcompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	//userId 是放到store里的
	//使用cookie
	//store := cookie.NewStore([]byte("secret"))                 ///设置session的密钥

	//使用redis
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "", //最大空闲连接数,tcp 连接地址,密码, key 和 value 的加密密钥
	//	[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	//	[]byte("0Pf2r0wZBpXVXLQNdpwCXN4ncnlnZSc3"),
	//) // 创建redis store
	//if err != nil {
	//	panic(err)
	//}
	//redisClient := redis.NewClient(&redis.Options{  // 创建redis client
	//	Addr: config.Config.Redis.Addr,
	//})
	//session中间件
	//store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
	//	[]byte("0Pf2r0wZBpXVXLQNdpwCXN4ncnlnZSc3")) // 创建memstore store
	//mystore := &sqlx_store.Store{}  // 创建sqlx store
	//server.Use(sessions.Sessions("mysession", store)) //设置session中间件
	//server.Use(middleware.NewLoginMiddlewareBuilder().Build()) //登录中间件

	//jwt中间件
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("users/login_sms/code/send").
		IgnorePaths("users/login_sms").
		Build())
	//server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build()) //在多长时间内允许多少请求
	return server
}
func initRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler { //初始化用户服务
	ud := dao.NewUserDao(db)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)
	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u

}

func initDB() *gorm.DB { // 初始化数据库
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

func main() {
	db := initDB()            // 初始化数据库
	server := initWebServer() // 初始化web服务
	initUser(db, initRedis())
	server.Run("127.0.0.1:8080")

}
