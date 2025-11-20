package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initVieprV1()
	////创建grpc服务
	app := InitApp()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	//启动服务
	err := app.server.Serve()
	log.Println(err)

}

//	func mainV1() {
//		initVieprV1()
//		server := grpc2.NewServer()
//		intrSvc := InitGrpcServer()
//		intrv1.RegisterInteractiveServiceServer(server, intrSvc)
//		//监听8090
//		l, err := net.Listen("tcp", ":8090")
//		if err != nil {
//			panic(err)
//		}
//		err = server.Serve(l)
//		log.Println(err)
//
// }
func initVieprV1() {
	cfile := pflag.String("config",
		"webook/interactive/config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	//实时监控配置变更
	viper.WatchConfig()
	//只能告诉你问价变了,不能告诉你,文件的那些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		//比较好的设计,他会在in 里面告诉你变更前的数据,和变更后的数据
		//更好的设计师,他会直接告诉你差异
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetStringMap("db.dsn"))
	})
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}
