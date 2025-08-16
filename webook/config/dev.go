//go:build !k8s
// +build !k8s

package config

var Config = config{
	DB: DBConfig{
		//本地连接
		//"root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local",
		"root:root@tcp(localhost:13316)/webook",
	},
	Redis: ReidsConfig{
		"locallhost:6379",
	},
}
