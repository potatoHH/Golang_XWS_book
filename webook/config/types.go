// 没有k8s这个标签
package config

type config struct {
	DB    DBConfig
	Redis ReidsConfig
}

type DBConfig struct {
	DSN string
}

type ReidsConfig struct {
	Addr string
}

type MysqlConfig struct {
}
