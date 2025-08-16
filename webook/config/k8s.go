//go:build k8s

// 使用k8s标签
package config

var Config = config{
	DB: DBConfig{
		"root:root@tcp(webook-live-mysql:11309)/webook",
	},
	Redis: ReidsConfig{
		"webook-live-redis:10379",
	},
}
