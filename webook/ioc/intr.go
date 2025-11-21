package ioc

import (
	intrv1 "Book_Exp/webook/api/proto/gen/intr/v1"
	"Book_Exp/webook/interactive/service"
	"Book_Exp/webook/internal/web/client"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGRPCClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr   string
		Secure bool
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if cfg.Secure {
		//上面要去加载你启用的证书之类的
		//启用https
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	
	cc, err := grpc.NewClient(cfg.Addr, opts...)
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewInteractiverServiceApadter(svc)
	return client.NewGreyScaleInteractiveServiceClient(local, remote)
}
