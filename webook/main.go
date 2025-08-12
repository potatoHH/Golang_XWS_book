package main

import (
	"Book_Exp/webook/internal/repository"
	"Book_Exp/webook/internal/repository/dao"
	"Book_Exp/webook/internal/service"
	"Book_Exp/webook/internal/web"
	"Book_Exp/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
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
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "", //最大空闲连接数,tcp 连接地址,密码, key 和 value 的加密密钥
		[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"),
		[]byte("0Pf2r0wZBpXVXLQNdpwCXN4ncnlnZSc3"),
	) // 创建redis store
	if err != nil {
		panic(err)
	}
	//mystore := &sqlx_store.Store{}  // 创建sqlx store
	server.Use(sessions.Sessions("mysession", store)) //设置session中间件
	//server.Use(middleware.NewLoginMiddlewareBuilder().Build()) //登录中间件

	//jwt中间件
	server.Use(middleware.NewLoginJwtMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		Build())
	server.Use()
	return server
}

func initUser(server *gin.Engine, db *gorm.DB) *web.UserHandler { //初始化用户服务
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewService(repo)
	u := web.NewUserHandler(svc)
	u.RegisterRoutes(server)
	return u

}

func initDB() *gorm.DB { // 初始化数据库
	dsn := "root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local" // 数据库连接
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
	initUser(server, db)      // 初始化用户
	server.Run("127.0.0.1:8080")

}
