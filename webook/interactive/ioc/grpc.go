package ioc

import (
	grpc2 "Book_Exp/webook/interactive/grpc"
	"Book_Exp/webook/pkg/grpcx"
	"Book_Exp/webook/pkg/logger"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(l logger.LoggerV1, intrServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Port      int      `yaml:"port"`
		EtcdAddrs []string `yaml:"etcdAddrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	//注册
	server := grpc.NewServer()
	intrServer.Register(server)
	return &grpcx.Server{
		Server:   grpc.NewServer(),
		Port:     cfg.Port,
		EtcdAddr: cfg.EtcdAddrs,
		Name:     "interactive",
		L:        l,
	}

}
