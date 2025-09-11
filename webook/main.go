package main

import (
	"github.com/spf13/viper"
)

func main() {
	initViperV1()
	server := InitWebServer()
	server.Run("127.0.0.1:8080")

}

//func initViper() {
//	viper.SetConfigFile("config/dev.yaml")
//	err := viper.ReadInConfig()
//	if err != nil {
//		panic(err)
//	}
//}

func initViperV1() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

//func initViperV2() {
//	cfile := pflag.String("config", "config/dev.yaml", "指定配置文件路径")
//	pflag.Parse() //解析命令行参数
//	viper.SetConfigFile(*cfile)
//}
