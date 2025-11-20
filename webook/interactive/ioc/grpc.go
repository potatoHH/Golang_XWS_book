package ioc

import (
	grpc2 "Book_Exp/webook/interactive/grpc"
	"Book_Exp/webook/pkg/grpcx"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(intrServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
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
		Server: grpc.NewServer(),
		Add:    cfg.Addr,
	}

}
